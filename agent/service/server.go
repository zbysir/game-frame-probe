package service

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bysir-zl/game-frame-probe/common/service"
	"sync"
	"fmt"
	"time"
	"github.com/bysir-zl/hubs/core/hubs"
	"github.com/bysir-zl/hubs/core/net/listener"
	"github.com/AsynkronIT/protoactor-go/remote"
	"strings"
	"github.com/bysir-zl/bygo/log"
)

const TAG = "agent"

var (
	id   string = "agent/1"
	addr string = "127.0.0.1"
	port int    = 8080
)

type Server struct {
	Server *service.Server
	*actor.PID
}

type Servers map[string]*Server
type ServerGroups map[string]Servers

func (p ServerGroups) SelectServer(serverType string) (server *Server, ok bool) {
	lock.RLock()
	defer lock.RUnlock()
	if servers, has := p[serverType]; has {
		if len(servers) != 0 {
			for _, s := range servers {
				// todo 这里可能需要返回服务器Pid的子Pid去处理用户请求, 但是这里每次都新建一个PID? 想一想
				server = s
				ok = true
				return
			}
		}
	}
	return 
}

type ServerChange struct {
	server *service.Server
	change service.ServerChange
}

var (
	lock            = sync.RWMutex{}
	stdServerGroups = ServerGroups{}
)

// 
func GetServers(serverType string) (servers map[string]*Server, has bool) {
	lock.RLock()
	if ss, ok := stdServerGroups[serverType]; ok {
		lock.RUnlock()
		if len(ss) != 0 {
			return ss, true
		}
	} else {
		lock.RUnlock()
	}
	return
}

func OnServerChange(msg *ServerChange) {
	server := msg.server
	id := server.Id
	switch msg.change {
	case service.SC_Online:
		serverAddr := fmt.Sprintf("%s:%d", server.Address, server.Port)
		// 不要连接自己
		if serverAddr == fmt.Sprintf("%s:%d", addr, port) {
			break
		}

		pid := actor.NewPID(serverAddr, id)
		serverType := getServerTypeFromId(id)
		if servers, ok := stdServerGroups[serverType]; ok {
			servers[id] = &Server{PID: pid}
		} else {
			stdServerGroups[serverType] = Servers{
				id: {PID: pid},
			}
		}
		log.InfoT(TAG, "server %s is readied", id)

	case service.SC_Offline:
		serverType := getServerTypeFromId(id)
		if servers, ok := stdServerGroups[serverType]; ok {
			delete(servers, server.Id)
		}

		log.InfoT(TAG, "server %s offline", id)
	}
}

func getServerTypeFromId(id string) string {
	return strings.Split(id, "/")[0]
}

func Run() {
	router := stdRouter
	agentPid, err := serverNode(router)
	if err != nil {
		panic(err)
	}

	serverCli(agentPid, router)

	// 注册服务
	manager := service.NewManagerEtcd()
	lease, err := manager.RegisterService(&service.Server{
		Id:      id,
		Name:    id,
		Address: addr,
		Port:    port,
	})
	if err != nil {
		panic(err)
	}

	// 监听服务变化
	manager.WatchServer(func(server *service.Server, change service.ServerChange) {
		OnServerChange(&ServerChange{server: server, change: change})
		return
	})

	for range time.Tick(time.Second * 5) {
		manager.UpdateServerTTL(lease)
	}
}

func serverNode(router *Router) (agentPid *actor.PID, err error) {
	agentActor := NewAgentActor(router)
	nodeAddr := fmt.Sprintf("%s:%d", addr, port)

	remote.Start(nodeAddr)
	agentPid, err = actor.SpawnNamed(actor.FromInstance(agentActor), id)
	return
}

func serverCli(agentPid *actor.PID, router *Router) {
	addr := "127.0.0.1:8081"

	cliServer = hubs.New(addr, listener.NewWs(), &ClientHandler{
		agentPid: agentPid,
		router:   router,
	})
	log.Info("serverCli started on", addr)
	go func() {
		err := cliServer.Run()
		if err != nil {
			log.Error("serverCli err", err)
		}
	}()
}

var cliServer *hubs.Server

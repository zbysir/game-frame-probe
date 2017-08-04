package service

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bysir-zl/game-frame-probe/common/service"
	"sync"
	"fmt"
	"time"
	"github.com/bysir-zl/hubs/core/hubs"
	"github.com/bysir-zl/hubs/core/net/listener"
	"github.com/bysir-zl/game-frame-probe/agent/act"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/bysir-zl/hubs/core/net/conn_wrap"
	"strings"
	"github.com/bysir-zl/game-frame-probe/proto/pbgo"
	"encoding/json"
	"github.com/bysir-zl/bygo/log"
)
const TAG = "agent"

var (
	id   string = "agent-1"
	addr string = "127.0.0.1"
	port int    = 8080
)

type Server struct {
	Server *service.Server
	*actor.PID
}

type Servers map[string]*Server
type ServerGroups map[string]Servers

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

func OnServerChange(msg *ServerChange, ctx actor.Context) {
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

type ClientReq struct {
	Body []byte
}

type ClientRsp struct {
	ToServerType string
	Message      interface{}
}

type ClientContext struct {
	Request *ClientReq
}

type ClientHandler struct {
	agentPid *actor.PID
}

func (p *ClientHandler) Server(server *hubs.Server, conn conn_wrap.Interface) {

	clientPid := actor.Spawn(actor.FromInstance(&act.ClientActor{Conn: conn}))

	firstMsg := make(chan struct{})
	// 5s没消息就关闭连接
	go func() {
		select {
		case <-firstMsg:
		case <-time.After(5 * time.Second):
			conn.Close()
		}
	}()

	ctx:=ClientContext{}

	go func() {
		bs, err := conn.Read()
		close(firstMsg)
		if err != nil {
			return
		}
		ctx.Request = &ClientReq{Body:bs}
		
		// 上线
		agent.RequestFuture(&pbgo.ClientOnline{Uid: uid}, 5*time.Second).Result()
		defer agent.Request(&pbgo.ClientOffline{Uid: uid}, clientPid)
		for {
			bs, err := conn.Read()
			if err != nil {
				return
			}

			m := Proto{}
			json.Unmarshal(bs, &m)
			toServerType := ""
			switch m.Cmd {
			case 1:
				toServerType = "game"
			}

			agent.Request(&ClientReq{body: bs}, )
		}

	}()

	return
}

func getServerTypeFromId(id string) string {
	return strings.Split(id, "/")[0]
}

func Run() {
	agentActor := act.NewAgentActor()
	nodeAddr := fmt.Sprintf("%s:%d", addr, port)

	remote.Start(nodeAddr)
	agentPid, err := actor.SpawnNamed(actor.FromInstance(agentActor), id)
	if err != nil {
		panic(err)
	}

	serverCli(agentPid)

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
		agentPid.Tell(&ServerChange{server: server, change: change})
		return
	})

	for range time.Tick(time.Second * 5) {
		manager.UpdateServerTTL(lease)
	}
}

func serverCli(agentPid *actor.PID) {
	addr := "127.0.0.1:8081"

	cliServer = hubs.New(addr, listener.NewWs(), &ClientHandler{
		agentPid: agentPid,
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

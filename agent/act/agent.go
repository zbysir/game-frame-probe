package act

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bysir-zl/game-frame-probe/proto/pbgo"
	"reflect"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/bysir-zl/bygo/log"
	"github.com/bysir-zl/hubs/core/hubs"
	"github.com/bysir-zl/hubs/core/net/listener"
	"github.com/bysir-zl/hubs/core/net/conn_wrap"
	"sync"
	"fmt"
	"time"
	"github.com/bysir-zl/game-frame-probe/common/service"
	"strings"
	"github.com/AsynkronIT/protoactor-go/eventstream"
)

type Server struct {
	*actor.PID
}

type ServerGroup struct {
	Type    string
	Servers map[string]*Server // id:server
}

type ServerChange struct {
	server *service.Server
	change service.ServerChange
}

type AgentActor struct {
	serverGroups map[string]*ServerGroup // type:serverGroup
	sync.Mutex
}

const TAG = "agent"

func (p *AgentActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *ServerChange:
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
			if group, ok := p.serverGroups[serverType]; ok {
				group.Servers[id] = &Server{PID: pid}
			} else {
				p.serverGroups[serverType] = &ServerGroup{
					Type: serverType,
					Servers: map[string]*Server{
						id: {PID: pid},
					},
				}
			}
			log.InfoT(TAG, "server %s is conned", id)
			
		case service.SC_Offline:
			serverType := getServerTypeFromId(id)
			if group, ok := p.serverGroups[serverType]; ok {
				delete(group.Servers, server.Id)
			}

			log.InfoT(TAG, "server %s offline", id)
		}

	case *pbgo.AgentForwardToSvr:
		telled := false
		if servers, ok := p.serverGroups[msg.ServerType]; ok {
			if len(servers.Servers) == 0 {
				break
			}
			
			for _, s := range servers.Servers {
				s.Request(msg, context.Spawn())
				telled = true
				break
			}
		}

		if !telled {
			log.WarnT(TAG, "forward server type %s is not found", msg.ServerType)
		}
	case *actor.Terminated:
		log.InfoT(TAG, "Terminated", msg, msg.Who, msg.AddressTerminated)
	case *actor.Started:
		log.InfoT(TAG, "actor started")

	default:
		log.Info(reflect.TypeOf(msg).String())
	}
}

var (
	id   string = "agent-1"
	addr string = "127.0.0.1"
	port int    = 8080
)

func Run() {
	stdActor := NewAgentActor()

	agentPid, err := serverNode(stdActor)
	if err != nil {
		panic(err)
	}
	serverCli(stdActor)

	eventstream.Subscribe(func(evt interface{}) {
		log.InfoT("actor watch", evt)
	})

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
		log.DebugT("test", "watch ", server, change)
		agentPid.Tell(&ServerChange{server: server, change: change})
		return
	})

	for range time.Tick(time.Second * 5) {
		manager.UpdateServerTTL(lease)
	}
}

func NewAgentActor() *AgentActor {
	return &AgentActor{
		serverGroups: map[string]*ServerGroup{},
	}
}

func serverNode(agentActor *AgentActor) (pid *actor.PID, err error) {
	nodeAddr := fmt.Sprintf("%s:%d", addr, port)

	remote.Start(nodeAddr)
	props := actor.FromInstance(agentActor)
	pid, err = actor.SpawnNamed(props, id)
	return
}

type ClientHandler struct {
	agentActor *AgentActor
}

func (p *ClientHandler) Server(server *hubs.Server, conn conn_wrap.Interface) {
	uid := ""
	defer func() {
		if uid != "" {
			server.UnSubscribe(conn, uid)
		}
	}()
	// 新建一个agent处理用户请求
	// 每一个用户一个actor
	agentPid := actor.Spawn(actor.FromInstance(p.agentActor))
	RunClientActor(conn, agentPid)
}

func getServerTypeFromId(id string) string {
	return strings.Split(id, "-")[0]
}

func serverCli(agentActor *AgentActor) {
	addr := "127.0.0.1:8081"

	cliServer = hubs.New(addr, listener.NewWs(), &ClientHandler{
		agentActor: agentActor,
	})
	log.InfoT(TAG, "serverCli started on", addr)
	go func() {
		err := cliServer.Run()
		if err != nil {
			log.ErrorT(TAG, "serverCli err", err)
		}
	}()
}

var cliServer *hubs.Server

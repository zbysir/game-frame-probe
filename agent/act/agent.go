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
	"encoding/json"
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
		switch msg.change {
		case service.SC_Online:
			addr := fmt.Sprintf("%s:%d", server.Address, server.Port)
			// 不要连接自己
			if addr == fmt.Sprintf("%s:%d", addr, port) {
				break
			}
			// 请求连接
			pid := actor.NewPID(addr, server.Id)
			_, err := pid.RequestFuture(&pbgo.AgentConnectReq{Agent: context.Self()}, 3*time.Second).Result()
			if err != nil {
				log.ErrorT(TAG, "conn server err:", err)
				break
			}

			id := server.Id
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

			context.Watch(pid)
		case service.SC_Offline:
			serverType := getServerTypeFromId(id)
			if group, ok := p.serverGroups[serverType]; ok {
				delete(group.Servers, server.Id)
			}
			log.InfoT(TAG, "server %s offline", id)
		}

	case *pbgo.AgentForwardToSvr:
		if servers, ok := p.serverGroups[msg.ServerType]; ok {
			if len(servers.Servers) == 0 {
				break
			}
			for _, s := range servers.Servers {
				s.Tell(msg)
				break
			}
		}
	case *pbgo.AgentForwardToCli:
		cliServer.SendToTopic(msg.Uid, msg.Body)
	case *actor.Terminated:
		log.InfoT(TAG, "Terminated", msg, msg.Who, msg.AddressTerminated)
	case *actor.Started:

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
	pid := actor.Spawn(actor.FromInstance(p.agentActor))
	for {
		bs, err := conn.Read()
		if err != nil {
			return
		}
		type Proto struct {
			Cmd  int
			Body string
		}
		p := Proto{}
		json.Unmarshal(bs, &p)
		switch p.Cmd {
		case 0:
			uid = p.Body
			server.Subscribe(conn, uid)
		case 1:
			msg := pbgo.AgentForwardToSvr{ServerType: "game", Body: []byte(p.Body), Uid: uid}
			pid.Tell(&msg)
		case 2:
			msg := pbgo.AgentForwardToSvr{ServerType: "room", Body: []byte(p.Body), Uid: uid}
			pid.Tell(&msg)
		}
	}
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

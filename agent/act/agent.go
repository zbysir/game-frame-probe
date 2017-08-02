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
)

type Server struct {
	*actor.PID
}

type ServerGroup struct {
	Type    string
	Servers map[string]*Server // id:server
}

type AgentActor struct {
	serverGroups map[string]*ServerGroup // type:serverGroup
	sync.Mutex
}

func (p *AgentActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
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
		log.InfoT("agent", "Terminated", msg.Who, msg.AddressTerminated)
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
	agentPid, err := serverNode()
	if err != nil {
		panic(err)
	}
	serverCli()

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

	manager.WatchServer(func(server *service.Server, change service.ServerChange) {
		switch change {
		case service.SC_Online:
			addr := fmt.Sprintf("%s:%d", server.Address, server.Port)
			pid := actor.NewPID(addr, server.Id)
			// 不要连接自己
			if addr == agentPid.Address {
				break
			}
			serverType := getServerTypeFromId(server.Id)
			if group, ok := StdActor.serverGroups[serverType]; ok {
				group.Servers[server.Id] = &Server{PID: pid}
			} else {
				StdActor.serverGroups[serverType] = &ServerGroup{
					Type: serverType,
					Servers: map[string]*Server{
						server.Id: {PID: pid},
					},
				}
			}
			pid.Tell(&pbgo.AgentConnect{Sender: agentPid})
		case service.SC_Offline:
			serverType := getServerTypeFromId(server.Id)
			if group, ok := StdActor.serverGroups[serverType]; ok {
				delete(group.Servers, server.Id)
			}
		}

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

func serverNode() (pid *actor.PID, err error) {
	nodeAddr := fmt.Sprintf("%s:%d", addr, port)

	remote.Start(nodeAddr)
	props := actor.FromInstance(StdActor)
	pid, err = actor.SpawnNamed(props, id)
	return
}

func clientHandle(server *hubs.Server, conn conn_wrap.Interface) {
	uid := ""
	defer func() {
		if uid != "" {
			server.UnSubscribe(conn, uid)
		}
	}()
	props := actor.FromInstance(StdActor)
	pid := actor.Spawn(props)
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

func serverCli() {
	addr := "127.0.0.1:8081"

	cliServer = hubs.New(addr, listener.NewWs(), clientHandle)
	log.InfoT("agent", "serverCli started on", addr)
	go func() {
		err := cliServer.Run()
		if err != nil {
			log.ErrorT("agent", "serverCli err", err)
		}
	}()
}

var cliServer *hubs.Server
var StdActor = NewAgentActor()

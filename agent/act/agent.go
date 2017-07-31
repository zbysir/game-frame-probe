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
)

type AgentActor struct {
	serverMap map[string]*actor.PID
	sync.Mutex
}

func (p *AgentActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *pbgo.AgentConnect:
		// 当有节点连接上这个网关
		p.Lock()
		p.serverMap[msg.ServerName] = msg.Sender
		p.Unlock()
		log.InfoT("Agent", msg.ServerName+" conned")

		props := actor.FromInstance(p)
		pid := actor.Spawn(props)
		msg.Sender.Tell(&pbgo.AgentConnected{Message: "Welcome!", Server: pid})

		context.Watch(msg.Sender)
	case *pbgo.AgentForwardToSvr:
		if pid, ok := p.serverMap[msg.ServerName]; ok {
			pid.Tell(msg)
		}
	case *pbgo.AgentForwardToCli:
		cliServer.SendToTopic(msg.Uid, msg.Body)
	case *actor.Terminated:
		context.Unwatch(msg.Who)
		log.InfoT("agent", "Terminated", msg.Who, msg.AddressTerminated)
	case *actor.Started:

	default:
		log.Info(reflect.TypeOf(msg).String())
	}
}

func Server() {
	_, err := serverNode()
	if err != nil {
		panic(err)
	}

	serverCli()
}

func NewAgentActor() *AgentActor {
	return &AgentActor{
		serverMap: map[string]*actor.PID{},
	}
}

func serverNode() (pid *actor.PID, err error) {
	nodeAddr := "127.0.0.1:8080"

	remote.Start(nodeAddr)
	props := actor.FromInstance(StdActor)
	pid, err = actor.SpawnNamed(props, "agent")
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
			msg := pbgo.AgentForwardToSvr{ServerName: "game", Body: []byte(p.Body), Uid: uid}
			pid.Tell(&msg)
		case 2:
			msg := pbgo.AgentForwardToSvr{ServerName: "room", Body: []byte(p.Body), Uid: uid}
			pid.Tell(&msg)
		}
	}
}

func serverCli() {
	addr := "127.0.0.1:8081"

	cliServer = hubs.New(addr, listener.NewWs(), clientHandle)
	cliServer.Run()
}

var cliServer *hubs.Server
var StdActor = NewAgentActor()

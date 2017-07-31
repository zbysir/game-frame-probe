package act

import (
	"github.com/bysir-zl/game-frame-probe/proto/pbgo"
	"reflect"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bysir-zl/bygo/log"
)

type GameActor struct {
	agent *actor.PID
}

func (p *GameActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *pbgo.AgentConnected:
		p.agent = msg.Server
	case *pbgo.AgentForwardToSvr:
		if msg.Uid != "" {
			p.agent.Tell(&pbgo.AgentForwardToCli{
				Uid:  msg.Uid,
				Body: []byte(`{"data":"hello"}`),
			})
		}
	default:
		log.Info(reflect.TypeOf(msg).String())
	}
}

func NewGameActor() *GameActor {
	return &GameActor{
	}
}

func Server() {
	pid, err := serverNode()
	if err != nil {
		panic(err)
	}
	connAgent(pid)

	<-(chan int)(nil)
}

func connAgent(pid *actor.PID) {
	agentAddr := "127.0.0.1:8080"
	aPid := actor.NewPID(agentAddr, "agent")
	aPid.Tell(&pbgo.AgentConnect{ServerName: "game", Sender: pid})
}

func serverNode() (pid *actor.PID, err error) {
	nodeAddr := "127.0.0.1:8090"

	remote.Start(nodeAddr)
	props := actor.FromInstance(NewGameActor())
	pid, err = actor.SpawnNamed(props, "game")
	return
}

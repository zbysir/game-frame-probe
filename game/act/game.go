package act

import (
	"github.com/bysir-zl/game-frame-probe/proto/pbgo"
	"reflect"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bysir-zl/bygo/log"
	"github.com/bysir-zl/game-frame-probe/common/app"
	"time"
	"fmt"
)

type GameActor struct {
	agent *actor.PID
}

func (p *GameActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *pbgo.AgentConnect:
		log.InfoT("Agent", msg.ServerName+" conned")
		p.agent = msg.Sender
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

var mPid *actor.PID

var (
	id   string = "game-1"
	addr string = "127.0.0.1"
	port int    = 8090
)

func Server() {
	pid, err := serverNode()
	if err != nil {
		panic(err)
	}
	mPid = pid

	app.Init()
	app.RegisterService(&app.Service{
		Id:      id,
		Name:    id,
		Address: addr,
		Port:    port,
	}, "10s")
	app.UpdateServerTTL(id, "pass")

	for range time.Tick(time.Second * 5) {
		app.UpdateServerTTL(id, "pass")
	}
}

func serverNode() (pid *actor.PID, err error) {
	nodeAddr := fmt.Sprintf("%s:%d", addr, port)

	remote.Start(nodeAddr)
	props := actor.FromInstance(NewGameActor())
	pid, err = actor.SpawnNamed(props, id)
	return
}

package act

import (
	"github.com/bysir-zl/game-frame-probe/proto/pbgo"
	"reflect"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bysir-zl/bygo/log"
	"time"
	"fmt"
	"github.com/bysir-zl/game-frame-probe/common/service"
)

const TAG = "game"

type GameActor struct {
}

func (p *GameActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *pbgo.ClientMessageReq:
		context.Respond(&pbgo.ClientMessageRsp{
			Body: []byte(fmt.Sprintf(`{"data":"hello %s"}`, string(msg.Uid))),
		})
	case *actor.Started:

	default:
		log.Info(reflect.TypeOf(msg).String())
	}
}

func NewGameActor() *GameActor {
	return &GameActor{
	}
}

var (
	id   string = "game/1"
	addr string = "127.0.0.1"
	port int    = 8090
)

func Server() {
	_, err := serverNode()
	if err != nil {
		panic(err)
	}

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

	for range time.Tick(time.Second * 5) {
		manager.UpdateServerTTL(lease)
	}
}

func serverNode() (pid *actor.PID, err error) {
	nodeAddr := fmt.Sprintf("%s:%d", addr, port)

	remote.Start(nodeAddr)
	props := actor.FromInstance(NewGameActor())
	pid, err = actor.SpawnNamed(props, id)
	return
}

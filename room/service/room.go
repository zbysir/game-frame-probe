package service

import (
	"github.com/bysir-zl/game-frame-probe/proto/pbgo"
	"reflect"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bysir-zl/bygo/log"
	"fmt"
)

type RoomActor struct {
}

func (p *RoomActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *pbgo.ClientMessageReq:
		ctx.Respond(&pbgo.ClientMessageRsp{
			Body: []byte(fmt.Sprintf(`{"data":"hello i am from room, %s"}`, string(msg.Uid))),
		})
	case *pbgo.GetServerClientActorReq:
		pid, err := actor.SpawnNamed(actor.FromInstance(NewRoomActor()), TAG+"1")
		if err != nil {
			log.ErrorT(TAG,"err ",err)
		}
		ctx.Respond(&pbgo.GetServerClientActorRsp{
			Pid: pid,
		})
	case *actor.Started:

	default:
		log.Info(reflect.TypeOf(msg).String())
	}
}

func NewRoomActor() *RoomActor {
	return &RoomActor{
	}
}

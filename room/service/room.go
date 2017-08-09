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

func (p *RoomActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *pbgo.ClientMessageReq:
		context.Respond(&pbgo.ClientMessageRsp{
			Body: []byte(fmt.Sprintf(`{"data":"hello i am from room, %s"}`, string(msg.Uid))),
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

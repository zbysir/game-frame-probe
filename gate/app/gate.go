package app

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bysir-zl/game-frame-probe/common/pbgo"
	"reflect"
	"github.com/bysir-zl/bygo/log"
)

type AgentActor struct {
	router *Router
}

func (p *AgentActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *pbgo.AgentForwardToSrv:
		// 服务器之间转发消息
		serverPid, message, err := p.router.RouteServer(msg, stdServerGroups)
		if err != nil {
			log.ErrorT(TAG, err)
			break
		}

		serverPid.Request(message, ctx.Sender())
	case *actor.Terminated:
		log.InfoT(TAG, "Terminated", msg, msg.Who, msg.AddressTerminated)
	case *actor.Started:
	case *pbgo.GatePing:
		ctx.Respond(&pbgo.GatePong{})

	default:
		log.Info(reflect.TypeOf(msg).String())
	}
}

func NewAgentActor(router *Router) *AgentActor {
	return &AgentActor{router: router}
}

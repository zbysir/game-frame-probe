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

	"fmt"


	"strings"
	service2 "github.com/bysir-zl/game-frame-probe/agent/service"
)


type AgentActor struct {
}

const TAG = "agent"

func (p *AgentActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *ClientRsp:
		p.OnClientMessage(msg, ctx)
	case *pbgo.ClientOnline:
		// todo 记录用户id->pid
	case *pbgo.ClientOffline:

	case *pbgo.AgentForwardToSvr:
		telled := false
		if servers, ok := p.serverGroups[msg.ServerType]; ok {
			if len(servers.Servers) == 0 {
				break
			}

			for _, s := range servers.Servers {
				s.Request(msg, ctx.Sender())
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

	default:
		log.Info(reflect.TypeOf(msg).String())
	}
}

func (p *AgentActor) OnClientMessage(msg *ClientRsp, ctx actor.Context) {
	telled := false
	if servers, ok := p.serverGroups[msg.ToServerType]; ok {
		if len(servers.Servers) != 0 {
			for _, s := range servers.Servers {

				s.Request(msg.Message, ctx.Sender())
				telled = true
			}
		}
	}

	if !telled {
		log.WarnT(TAG, "forward server type %s is not found", msg.ToServerType)
	}
}


func NewAgentActor() *AgentActor {
	return &AgentActor{	}
}

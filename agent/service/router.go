package service

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"errors"
	"github.com/bysir-zl/game-frame-probe/proto/pbgo"
	"github.com/bysir-zl/bygo/log"
	"reflect"
)

var stdRouter = &Router{}

type Router struct {
}

func (p *Router) Route(ctx *SenderContext, servers ServerGroups) (serverPid *actor.PID, message interface{}, err error) {
	serverType := ""

	switch msg := ctx.Request.(type) {
	case *ClientReq:
		// 用户请求转发
		message = &pbgo.ClientMessageReq{Body: msg.Body, Uid: "sbsb2"}
		serverType = "game"
	default:
		log.InfoT("test", reflect.TypeOf(msg))
		// 其他的直接转发
		message = ctx.Request
		serverType = "game"
	}

	server, ok := servers.SelectServer(serverType)
	if !ok {
		err = errors.New("404")
		return
	}
	serverPid = server.PID

	return
}

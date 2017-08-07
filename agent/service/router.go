package service

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"errors"
	"github.com/bysir-zl/game-frame-probe/proto/pbgo"
)

var stdRouter = &Router{}

type Router struct {
}

// 请求路由
// 用户请求
func (p *Router) RouteClient(ctx *ClientContext, servers ServerGroups) (serverPid *actor.PID, message interface{}, err error) {
	serverType := ""
	uid := "sbsb2"

	message = &pbgo.ClientMessageReq{Body: ctx.Request.Body, Uid: uid}
	serverType = "game"

	server, ok := servers.SelectServer(serverType)
	if !ok {
		err = errors.New("404")
		return
	}
	serverPid = server.PID

	return
}

// 请求路由
// 服务器间的转发
func (p *Router) RouteServer(req *pbgo.AgentForwardToSvr, servers ServerGroups) (serverPid *actor.PID, message interface{}, err error) {
	server, ok := servers.SelectServer(req.ServerType)
	if !ok {
		err = errors.New("404")
		return
	}
	serverPid = server.PID
	message = req
	return
}

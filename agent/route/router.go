package route

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bysir-zl/game-frame-probe/agent/service"
)
//
//type ContextI interface {
//	
//}
//
//type Context struct {
//	Uid  string
//	req *act.ClientReq
//}
//
//func (p *Context) Req ()(req *ClientReq) {
//	return p.req
//}
//func (p *Context) SetReq (req *ClientReq)() {
//	p.req =req
//}
//
//type ClientHandler struct {
//}
//
//
//func (p *ClientHandler) Message(req *ClientReq,context interface{},) {
//
//}
//func (p *ClientHandler) DisConn() {
//
//}

type Router struct {
}

func (p *Router) Route(ctx *ClientContext, servers service.ServerGroups) (serverPid *actor.PID, message interface{}) {
	
	return "game"
}

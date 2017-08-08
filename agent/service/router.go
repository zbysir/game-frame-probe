package service

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"errors"
	"github.com/bysir-zl/game-frame-probe/proto/pbgo"
	"encoding/json"
)

var stdRouter = &Router{}

type Router struct {
}

type ClientMessage struct {
	Cmd  int `json:"cmd"`
	Body string `json:"body"`
}

const (
	Cmd_Login = iota + 1
	Cmd_Game  
	Cmd_Room  
)

// 请求路由
// 用户请求
func (p *Router) RouteClient(ctx *ClientContext, servers ServerGroups) (serverPid *actor.PID, message interface{}, err error) {
	serverType := ""

	clientMessage := ClientMessage{}
	json.Unmarshal(ctx.Request.Body, &clientMessage)
	switch clientMessage.Cmd {
	case Cmd_Login:
		// 用户登录
		ctx.SetValue("uid", clientMessage.Body)
		ctx.ClientPid.Tell(&pbgo.ClientMessageRsp{
			Body: []byte(`{"cmd":1,"body":"200"}`),
		})
		return
		
	case Cmd_Game:
		uid, ok := ctx.GetValue("uid")
		if !ok {
			ctx.ClientPid.Tell(&pbgo.ClientCloseRsq{
				Body: []byte(`{"cmd":1,"body":"bad message"}`),
			})
			return
		}
		message = &pbgo.ClientMessageReq{Body: ctx.Request.Body, Uid: uid.(string)}
		serverType = "game"
	case Cmd_Room:
		uid, ok := ctx.GetValue("uid")
		if !ok {
			ctx.ClientPid.Tell(&pbgo.ClientCloseRsq{
				Body: []byte(`{"cmd":1,"body":"bad message"}`),
			})
			return
		}
		message = &pbgo.ClientMessageReq{Body: ctx.Request.Body, Uid: uid.(string)}
		serverType = "room"
	}

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
func (p *Router) RouteServer(req *pbgo.AgentForwardToSrv, servers ServerGroups) (serverPid *actor.PID, message interface{}, err error) {
	server, ok := servers.SelectServer(req.ServerType)
	if !ok {
		err = errors.New("404")
		return
	}
	serverPid = server.PID
	message = req
	return
}

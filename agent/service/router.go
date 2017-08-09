package service

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"errors"
	"github.com/bysir-zl/game-frame-probe/proto/pbgo"
	"github.com/bysir-zl/game-frame-probe/common"
	"github.com/bysir-zl/game-frame-probe/common/client_msg"
	"github.com/bysir-zl/bygo/log"
	"github.com/gogo/protobuf/proto"
)

var stdRouter = &Router{}

type Router struct {
}

// 获取消息应当转发给的服务器
func (p *Router) GetClientMessageServerType(ctx *ClientContext) (serverType string, uid string, shouldForward bool) {
	serverType = ""

	clientMessage := client_msg.GetProto(ctx.Request.Body)
	uidX, ok := ctx.GetValue("uid")
	if ok {
		uid = uidX.(string)
	}

	switch clientMessage.Cmd {
	case common.CMD_Login:
		// 用户登录
		ctx.SetValue("uid", clientMessage.Body)
		ctx.ClientPid.Tell(&pbgo.ClientMessageRsp{
			Body: []byte(`{"cmd":1,"body":"200"}`),
		})
		return

	case common.CMD_JoinRoom, common.CMD_InitPlayer, common.CMD_Move:
		// 应该转发到game
		if uidX == "" {
			ctx.ClientPid.Tell(&pbgo.ClientCloseRsq{
				Body: []byte(`{"cmd":1,"body":"bad message"}`),
			})
			return
		}
		serverType = "game"
	}

	shouldForward = true
	return
}

// 请求路由
// 用户请求
func (p *Router) RouteClient(ctx *ClientContext, servers ServerGroups) {
	serverType := ""
	var message proto.Message

	clientMessage := client_msg.GetProto(ctx.Request.Body)
	switch clientMessage.Cmd {
	case common.CMD_Login:
		// 用户登录
		ctx.SetValue("uid", clientMessage.Body)
		ctx.ClientPid.Tell(&pbgo.ClientMessageRsp{
			Body: []byte(`{"cmd":1,"body":"200"}`),
		})
		return

	case common.CMD_JoinRoom, common.CMD_InitPlayer, common.CMD_Move:
		uid, ok := ctx.GetValue("uid")
		if !ok {
			ctx.ClientPid.Tell(&pbgo.ClientCloseRsq{
				Body: []byte(`{"cmd":1,"body":"bad message"}`),
			})
			return
		}
		message = &pbgo.ClientMessageReq{Body: ctx.Request.Body, Uid: uid.(string)}
		serverType = "game"
	case 10086:
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
		log.ErrorT(TAG, "not found game type of "+serverType)
		return
	}
	server.PID.Request(message, ctx.ClientPid)

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

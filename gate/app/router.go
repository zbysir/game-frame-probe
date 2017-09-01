package app

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"errors"
	"github.com/bysir-zl/game-frame-probe/common/pbgo"
	"github.com/bysir-zl/game-frame-probe/common"
	"github.com/bysir-zl/game-frame-probe/common/client_msg"
)

var stdRouter = &Router{}

type Router struct {
}

// 获取消息应当转发给的服务器
func (p *Router) RouteClient(ctx *ClientContext) (serverType string, shouldForward bool) {
	serverType = ""

	clientMessage := client_msg.GetProto(ctx.Request.Body)

	switch clientMessage.Cmd {
	case common.CMD_Login:
		// 用户登录
		ctx.Uid = clientMessage.Body
		ctx.SetValue("uid", clientMessage.Body)
		ctx.Pid.Tell(&pbgo.ClientMessageRsp{
			Body: []byte(`{"cmd":1,"body":"200"}`),
		})
		return

	case common.CMD_JoinRoom, common.CMD_InitPlayer, common.CMD_Broad:
		// 应该转发到game
		if ctx.Uid == "" {
			ctx.Pid.Tell(&pbgo.ClientCloseRsq{
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

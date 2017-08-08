package service

import (
	"github.com/bysir-zl/hubs/core/net/conn_wrap"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bysir-zl/game-frame-probe/proto/pbgo"
	"github.com/bysir-zl/bygo/log"
	"github.com/bysir-zl/hubs/core/hubs"
	"time"
	"sync"
)

// 这里实现 client -> agent -> game 之间通信
// client -tcp-AgentForwardToSvr> agent -AgentForwardToSvr> game -ClientMessageRsp> client

type ClientActor struct {
	// 网络连接
	Conn conn_wrap.Interface
}

type ClientConnReq struct {
}

type ClientConnRsp struct {
	AgentPid *actor.PID
}

type ClientClose struct {
}

func (p *ClientActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *pbgo.ClientMessageRsp:
		err := p.Conn.Write(msg.Body)
		if err != nil {
			log.ErrorT("c", err)
		}
	case *pbgo.ClientCloseRsq:
		p.Conn.WriteSync(msg.Body)
		p.Conn.Close()
	}
}

type ClientReq struct {
	Body []byte
}

type ClientContext struct {
	Request   *ClientReq // 客服端请求
	ClientPid *actor.PID // 实现直接返回消息
	data      map[string]interface{}
	l         *sync.RWMutex
}

func (p *ClientContext) SetValue(key string, value interface{}) {
	if p.data == nil {
		p.data = map[string]interface{}{}
		p.l = &sync.RWMutex{}
	}
	p.l.Lock()
	p.data[key] = value
	p.l.Unlock()
}

func (p *ClientContext) GetValue(key string) (value interface{}, ok bool) {
	if p.data == nil {
		return
	}

	p.l.RLock()
	value, ok = p.data[key]
	p.l.RUnlock()
	return
}

type ClientHandler struct {
	agentPid *actor.PID
	router   *Router
}

func (p *ClientHandler) Server(server *hubs.Server, conn conn_wrap.Interface) {
	clientPid := actor.Spawn(actor.FromInstance(&ClientActor{Conn: conn}))

	firstMsg := make(chan struct{})
	// 5s没消息就关闭连接
	go func() {
		select {
		case <-firstMsg:
		case <-time.After(5 * time.Second):
			clientPid.Tell(&pbgo.ClientCloseRsq{Body: []byte("timeout of auth")})
		}
	}()

	// 一个请求一个上下文, 用来存储登录信息等
	ctx := &ClientContext{
		ClientPid: clientPid,
	}

	go func() {
		// 读一次, 用来实现超时没消息关闭
		bs, err := conn.Read()
		close(firstMsg)
		if err != nil {
			return
		}
		ctx.Request = &ClientReq{Body: bs}
		serverPid, message, err := p.router.RouteClient(ctx, stdServerGroups)
		if err != nil {
			log.ErrorT(TAG, err)
		} else if serverPid != nil {
			serverPid.Request(message, clientPid)
		}

		for {
			bs, err := conn.Read()
			if err != nil {
				return
			}
			ctx.Request = &ClientReq{Body: bs}
			serverPid, message, err := p.router.RouteClient(ctx, stdServerGroups)
			if err != nil {
				log.ErrorT(TAG, err)
				continue
			} else if serverPid != nil {
				serverPid.Request(message, clientPid)
			}
		}

	}()

	return
}

package app

import (
	"github.com/bysir-zl/hubs/core/net/conn_wrap"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bysir-zl/game-frame-probe/common/pbgo"
	"github.com/bysir-zl/bygo/log"
	"github.com/bysir-zl/hubs/core/hubs"
	"time"
	"sync"
	"errors"
	"sync/atomic"
)

// 这里实现 client -> gate -> 游戏服务器节点 之间通信

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

func (p *ClientActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *pbgo.ClientMessageRsp:
		err := p.Conn.Write(msg.Body)
		if err != nil {
			log.ErrorT("c", err)
		}
	case *pbgo.ClientCloseRsq:
		p.Conn.WriteSync(msg.Body)
		p.Conn.Close()
		ctx.Self().Stop()
	}
}

type ClientReq struct {
	Body []byte
}

// gate 维持着客户端与服务器节点的连接, 当客户端上线, 消息, 下线, 都会通知响应的服务器节点; 一个客户端可以连接多个服务器节点;
type ClientContext struct {
	Uid          string             // 连接的唯一编号
	Request      *ClientReq         // 客服端请求
	Pid          *actor.PID         // 实现直接返回消息
	ConnedServer map[string]*Server // 连接过的服务器, serverType->server 
	data         map[string]interface{}
	l            *sync.RWMutex
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

func (p *ClientContext) DelValue(key string) {
	if p.data == nil {
		return
	}

	p.l.Lock()
	delete(p.data, key)
	p.l.Unlock()
	return
}

// 通知节点客户端断开连接
func (p *ClientContext) Close() {
	// 通知节点断开连接
	for _, s := range p.ConnedServer {
		s.Request(&pbgo.ClientDisconnectReq{Uid: p.Uid}, p.Pid)
	}

	// 通知关闭actor与连接
	p.Pid.Tell(&pbgo.ClientCloseRsq{Body: []byte("timeout of auth")})
}

// 向节点发消息
func (p *ClientContext) SendOrConnServer(serverType string, bs []byte) (err error) {
	if s, ok := p.ConnedServer[serverType]; ok {
		// 连接过, 就直接转发
		s.Request(&pbgo.ClientMessageReq{Uid: p.Uid, Body: bs}, p.Pid)
	} else {
		// 没连接过 就找到服务, 连接并转发
		s, ok := stdServerGroups.SelectServer(serverType)
		if !ok {
			return errors.New("bad cmd, can't found server, serverType:" + serverType)
		} else {
			// 存储服务器, 方便在Client断开的时候通知服务器
			p.ConnedServer[serverType] = s
			s.Request(&pbgo.ClientConnectReq{Uid: p.Uid}, p.Pid)
			s.Request(&pbgo.ClientMessageReq{Uid: p.Uid, Body: bs}, p.Pid)
		}
	}

	return
}

type ClientHandler struct {
	agentPid    *actor.PID
	router      *Router
	wg          sync.WaitGroup
	clientCount uint64 // 在线的客户端数量
}

func (p *ClientHandler) Server(server *hubs.Server, conn conn_wrap.Interface) {
	clientPid := actor.Spawn(actor.FromInstance(&ClientActor{Conn: conn}))

	firstMsg := make(chan struct{})

	// 一个请求一个上下文, 用来存储登录信息等
	ctx := &ClientContext{
		Pid:          clientPid,
		ConnedServer: map[string]*Server{},
	}

	// 5s没消息就关闭连接
	go func() {
		select {
		case <-firstMsg:
		case <-time.After(5 * time.Second):
			ctx.Close()
		}
	}()

	go func() {
		// 读一次, 用来实现超时没消息关闭
		bs, err := conn.Read()
		close(firstMsg)
		if err != nil {
			return
		}

		p.wg.Add(1)
		atomic.AddUint64(&p.clientCount, 1)
		// 关闭连接
		defer func() {
			ctx.Close()
			p.wg.Done()
			atomic.AddUint64(&p.clientCount, ^uint64(1 - 1))
		}()

		ctx.Request = &ClientReq{Body: bs}
	handle:
		serverType, should := p.router.RouteClient(ctx)
		// 一旦有了Uid, 那么Pid应该替换, Pid与Uid一一对应, 实现在客户端重连后在各个服务器保存的Pid还能生效;
		if ctx.Uid != "" && ctx.Uid != ctx.Pid.Id {
			ctx.Pid.Stop()
			ctx.Pid, err = actor.SpawnNamed(actor.FromInstance(&ClientActor{Conn: conn}), "U"+ctx.Uid)
			if err != nil {
				log.ErrorT(TAG, "SpawnNamed err:", err)
				return
			}
		}

		if should {
			err := ctx.SendOrConnServer(serverType, bs)
			if err != nil {
				ctx.Pid.Tell(&pbgo.ClientMessageRsp{Body: []byte(err.Error())})
			}
		}
		// end

		bs, err = conn.Read()
		if err != nil {
			return
		}
		ctx.Request = &ClientReq{Body: bs}

		goto handle
	}()

	return
}

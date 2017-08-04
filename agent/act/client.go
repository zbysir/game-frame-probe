package act

import (
	"github.com/bysir-zl/hubs/core/net/conn_wrap"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bysir-zl/game-frame-probe/proto/pbgo"
	"github.com/bysir-zl/bygo/log"
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
	case *ClientClose:
		p.Conn.Close()
	}
}

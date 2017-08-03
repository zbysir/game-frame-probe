package act

import (
	"github.com/bysir-zl/hubs/core/net/conn_wrap"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bysir-zl/game-frame-probe/proto/pbgo"
	"time"
	"encoding/json"
	"github.com/bysir-zl/bygo/log"
)

// 这里实现 client -> agent -> game 之间通信
// client -tcp-AgentForwardToSvr> agent -AgentForwardToSvr> game -ClientMessageRsp> client

type ClientActor struct {
	// 网络连接
	conn conn_wrap.Interface
}

type ClientClose struct {
}

type ClientMessage struct {
	uid          string
	toServerType string
	body         []byte
}

func (p *ClientActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *pbgo.ClientMessageRsp:
		err := p.conn.Write(msg.Body)
		if err != nil {
			log.ErrorT("c", err)
		}
	case *ClientClose:
		p.conn.Close()
	}
}

// 建立client与agent关联
func StartClientActorRecvice(conn conn_wrap.Interface, agent *actor.PID) {
	clientActor := &ClientActor{conn: conn}
	pid := actor.Spawn(actor.FromInstance(clientActor))

	firstMsg := make(chan struct{})
	// 5s没消息就关闭连接
	go func() {
		select {
		case <-firstMsg:
		case <-time.After(5 * time.Second):
			conn.Close()
		}
	}()

	go func() {
		bs, err := conn.Read()
		close(firstMsg)
		if err != nil {
			return
		}

		type Proto struct {
			Cmd  int
			Body string
		}
		m := Proto{}
		json.Unmarshal(bs, &m)
		uid := ""

		if m.Cmd != 0 {
			conn.Close()
			return
		}
		// 登陆
		uid = m.Body

		// 上线
		agent.Request(&pbgo.ClientOnline{Uid: uid}, pid)
		defer agent.Request(&pbgo.ClientOffline{Uid: uid}, pid)
		for {
			bs, err := conn.Read()
			if err != nil {
				return
			}

			m := Proto{}
			json.Unmarshal(bs, &m)
			toServerType := ""
			switch m.Cmd {
			case 1:
				toServerType = "game"
			}

			agent.Request(&ClientMessage{body: bs, toServerType: toServerType, uid: uid}, pid)
		}

	}()

	return
}

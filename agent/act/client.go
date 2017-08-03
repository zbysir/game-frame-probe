package act

import (
	"github.com/bysir-zl/hubs/core/net/conn_wrap"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bysir-zl/game-frame-probe/proto/pbgo"
	"encoding/json"
)

// 这里实现 client -> agent -> game 之间通信
// client -tcp-AgentForwardToSvr> agent -AgentForwardToSvr> game -ClientMessageRsp> client

type ClientActor struct {
	// 网络连接
	conn conn_wrap.Interface
}

func (p *ClientActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *pbgo.ClientMessageRsp:
		p.conn.Write(msg.Body)
	}
}

// 建立client与agent关联
func RunClientActor(conn conn_wrap.Interface, agent *actor.PID) *ClientActor {
	playerActor := &ClientActor{conn: conn}
	pid := actor.Spawn(actor.FromInstance(playerActor))

	go func() {
		for {
			bs, err := conn.Read()
			if err != nil {
				return
			}
			type Proto struct {
				Cmd  int
				Body string
			}
			p := Proto{}
			json.Unmarshal(bs, &p)
			switch p.Cmd {
			case 1:
				msg := pbgo.AgentForwardToSvr{ServerType: "game", Body: []byte(p.Body)}
				agent.Request(&msg, pid)
			}
		}
	}()

	return playerActor
}

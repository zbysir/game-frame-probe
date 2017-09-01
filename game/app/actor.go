package app

import (
	"github.com/bysir-zl/game-frame-probe/common/pbgo"
	"reflect"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bysir-zl/bygo/log"
	"github.com/bysir-zl/game-frame-probe/common/client_msg"
	"github.com/bysir-zl/game-frame-probe/common"
	"encoding/json"
	"time"
)

const frameDuration = time.Second / 10

// 一个房间一个actor
type GameActor struct {
	Players     map[string]*Player
	MessageList chan []byte
}

func (p *GameActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *pbgo.ClientConnectReq:
		log.InfoT(TAG, "user", msg.Uid, "connected")
	case *pbgo.ClientDisconnectReq:
		delete(p.Players, msg.Uid)
		log.InfoT(TAG, "user", msg.Uid, "disconnected")
	case *pbgo.ClientMessageReq:
		// 只要有消息就加入这个列表
		if _, ok := p.Players[msg.Uid]; !ok {
			p.Players[msg.Uid] = &Player{
				Uid: msg.Uid,
				PID: ctx.Sender(),
			}
		}

		pt := client_msg.GetProto(msg.Body)
		switch pt.Cmd {
		case common.CMD_JoinRoom:
			bs, _ := json.Marshal(p.Players)
			rsp := client_msg.NewProto(common.CMD_InitPlayer, bs)

			// 广播有人加入了房间
			p.ReadyBroadToPlayer(rsp)
		case common.CMD_Broad:
			player := Broad{
				Uid:  msg.Uid,
				Body: string(msg.Body),
			}
			bs, _ := json.Marshal(&player)
			rsp := client_msg.NewProto(common.CMD_Broad, bs)

			// 广播
			p.ReadyBroadToPlayer(rsp)
		}
	case *actor.Started:
		go func() {
			log.InfoT(TAG, "started")
			t := time.NewTicker(frameDuration)
			// 发送逻辑帧, 空帧也发
			for {
				select {
				case <-t.C:
					msgs := []string{}
					for {
						select {
						case msg := <-p.MessageList:
							msgs = append(msgs, string(msg))
						default:
							goto end
						}
					}
				end:
				//if len(msgs) == 0 {
				//	break
				//}
					for _, p := range p.Players {
						bs, _ := json.Marshal(&msgs)
						rsp := client_msg.NewProto(common.CMD_Logic, bs)
						p.Tell(&pbgo.ClientMessageRsp{
							Body: rsp,
						})
					}
				}
			}
			log.InfoT(TAG, "end")
		}()
	default:
		log.Info(reflect.TypeOf(msg).String())
	}
}

// 准备广播
func (p *GameActor) ReadyBroadToPlayer(msg []byte) {
	p.MessageList <- msg
}

func NewGameActor() *GameActor {
	return &GameActor{
		Players:     map[string]*Player{},
		MessageList: make(chan []byte, 1000),
	}
}

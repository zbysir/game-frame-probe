package service

import (
	"github.com/bysir-zl/game-frame-probe/proto/pbgo"
	"reflect"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bysir-zl/bygo/log"
	"github.com/bysir-zl/game-frame-probe/common/client_msg"
	"github.com/bysir-zl/game-frame-probe/common"
	"encoding/json"
	"time"
	"github.com/gogo/protobuf/proto"
)

type GameActorManager struct {
	ClientShouldJoin map[string]*actor.PID // 客户端应该进入的游戏actor
}

type GameActor struct {
	Players     map[string]*Player
	MessageList chan []byte
}

func (p *GameActorManager) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *pbgo.ClientMessageReq:
		// 分配给指定的actor
		if pid, ok := p.ClientShouldJoin[msg.Uid]; ok {
			pid.Tell(msg)
		} else {
			// 应当从数据库获取得到这个用户应该给那个游戏actor通信
			// 在pvp中, 在一起游戏(一个召唤师峡谷)的玩家应当在同一个actor.
			shouldJoinGameId := "1"
			pid, err := actor.SpawnNamed(actor.FromInstance(NewGameActor()), shouldJoinGameId)
			if err != nil {
				ctx.Respond(proto.ErrNil)
				return
			}
			p.ClientShouldJoin[msg.Uid] = pid
		}

	case *actor.Started:
	}
}

func NewGameActorManager() *GameActorManager {
	return &GameActorManager{
		ClientShouldJoin: map[string]*actor.PID{},
	}
}

func (p *GameActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *pbgo.ClientMessageReq:
		pt := client_msg.GetProto(msg.Body)
		switch pt.Cmd {
		case common.CMD_JoinRoom:
			player := Player{
				Uid:   msg.Uid,
				Point: Point{X: 50, Y: 50},
				PID:   ctx.Sender(),
			}
			p.Players[msg.Uid] = &player

			bs, _ := json.Marshal(p.Players)
			rsp := client_msg.NewProto(common.CMD_InitPlayer, bs)

			p.ReadyBroadToPlayer(rsp)
		case common.CMD_Broad:
			player := Broad{
				Uid:  msg.Uid,
				Body: string(msg.Body),
			}
			bs, _ := json.Marshal(&player)
			rsp := client_msg.NewProto(common.CMD_Broad, bs)

			p.ReadyBroadToPlayer(rsp)
		}
	case *actor.Started:
		go func() {
			log.InfoT(TAG, "started")
			t := time.NewTicker(time.Second / 20) // ,每20分之一秒发布一个逻辑帧
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

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
		}
	case int:
		pid, _ := actor.SpawnNamed(actor.FromInstance(NewGameActor()), TAG+"1")
		p.ClientShouldJoin["1"] = pid
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
				Point: Point{X: 1, Y: 1},
				PID:   ctx.Sender(),
			}
			p.Players[msg.Uid] = &player

			bs, _ := json.Marshal(&player)
			rsp := client_msg.NewProto(common.CMD_InitPlayer, bs)

			p.ReadyBroadToPlayer(rsp)
		case common.CMD_Move:
			m := client_msg.GetMove(msg.Body)
			player := Player{
				Uid:   msg.Uid,
				Angle: m.Angle,
				Speed: m.Speed,
			}
			bs, _ := json.Marshal(&player)
			rsp := client_msg.NewProto(common.CMD_Move, bs)

			p.ReadyBroadToPlayer(rsp)
		}
	case *actor.Started:
		go func() {
			log.InfoT(TAG, "started")
			t := time.NewTicker(time.Second / 20)
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
					if len(msgs) == 0 {
						break
					}
					for _, p := range p.Players {
						bs, _ := json.Marshal(&msgs)
						rsp := client_msg.NewProto(common.CMD_Logic, bs)
						log.InfoT("test", "rsp", string(rsp))
						if p.Uid == "2" {
							p.Tell(&pbgo.ClientMessageRsp{
								Body: rsp,
							})
						} else {
							p.Tell(&pbgo.ClientMessageRsp{
								Body: rsp,
							})
						}
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

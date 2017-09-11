package app

import (
	"encoding/json"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bysir-zl/bygo/log"
	"github.com/bysir-zl/game-frame-probe/common"
	"github.com/bysir-zl/game-frame-probe/common/client_msg"
	"github.com/bysir-zl/game-frame-probe/common/pbgo"
	"reflect"
	"time"
)

const frameDuration = time.Second / 10

// 一个房间一个actor
type GameActor struct {
	Id          string
	players     map[string]*Player
	messageList chan []byte
	stopC       chan struct{}
	manager     *actor.PID
	isPause     bool
	maxIdleTime int64 // 最大暂停时间, 超过就关闭这个actor
}

type ActorPauseMsg struct {
	IsPause bool
}

func (p *GameActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *pbgo.ClientConnectReq:
		if _, ok := p.players[msg.Uid]; !ok {
			p.players[msg.Uid] = &Player{
				Uid: msg.Uid,
				PID: ctx.Sender(),
			}
		}

		log.InfoT(TAG, "user", msg.Uid, "connected")
	case *pbgo.ClientDisconnectReq:
		delete(p.players, msg.Uid)
		if len(p.players) == 0 {
			ctx.Self().Tell(&ActorPauseMsg{IsPause: true})
		}
		log.InfoT(TAG, "user", msg.Uid, "disconnected")
	case *pbgo.ClientMessageReq:
		pt := client_msg.GetProto(msg.Body)
		switch pt.Cmd {
		case common.CMD_JoinRoom:
			bs, _ := json.Marshal(p.players)
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
			log.InfoT(TAG, "actor %s started", p.Id)
			defer func() {
				log.InfoT(TAG, "actor %s stopped", p.Id)
			}()
			t := time.NewTicker(frameDuration)
			var pauseStartTime int64 = 0
			// 发送逻辑帧, 空帧也发
			for {
				select {
				case <-p.stopC:
					return
				case <-t.C:
					if p.isPause {
						if pauseStartTime == 0 {
							pauseStartTime = time.Now().Unix()
						} else {
							// 超过最大暂停时间就关闭这个actor
							if time.Now().Unix()-pauseStartTime > p.maxIdleTime {
								ctx.Self().Stop()
								return
							}
						}
						continue
					}
					pauseStartTime = 0

					msgs := struct {
						Msg []string `json:"msg"`
					}{
						Msg: []string{},
					}

					for {
						select {
						case msg := <-p.messageList:
							msgs.Msg = append(msgs.Msg, string(msg))
						default:
							goto end
						}
					}
				end:
				//if len(msgs) == 0 {
				//	break
				//}
					for _, p := range p.players {
						bs, _ := json.Marshal(&msgs)
						rsp := client_msg.NewProto(common.CMD_Logic, bs)
						p.Tell(&pbgo.ClientMessageRsp{
							Body: rsp,
						})
					}
				}
			}
		}()
	case *ActorPauseMsg:
		p.isPause = msg.IsPause
	case *actor.Stopping:
	case *actor.Stopped:
		// 通知manager删除自己
		p.manager.Tell(&ActorStopMsg{ActorId: p.Id})
		// 关闭循环
		close(p.stopC)
	default:
		log.Info(reflect.TypeOf(msg).String())
	}
}

// 准备广播
func (p *GameActor) ReadyBroadToPlayer(msg []byte) {
	if p.isPause {
		return
	}
	p.messageList <- msg
}

func NewGameActor(manager *actor.PID, actorId string) *GameActor {
	return &GameActor{
		Id:          actorId,
		players:     map[string]*Player{},
		messageList: make(chan []byte, 1000),
		stopC:       make(chan struct{}, 1),
		manager:     manager,
		maxIdleTime: 60 * 3, // 最多暂停3分钟就关闭这个actor
	}
}

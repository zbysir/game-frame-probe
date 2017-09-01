package app

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bysir-zl/game-frame-probe/common/pbgo"
	"github.com/bysir-zl/bygo/log"
)

type GameActorManager struct {
	// 客户端应该进入哪个房间, 在连接的时候, 需要指定房间Id; 在之后的通信就不需要了; 为了down机恢复,需要持久化此字段
	ClientShouldJoin map[string]string     // uid=>roomId
	Actors           map[string]*actor.PID // 所有已经生成的actor, name=>pid
}

func (p *GameActorManager) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *pbgo.ClientConnectReq:
		// 在第一次客户端向这个节点发消息的时候会连接这个节点, 
		// 如果本节点重启后, 也不再有此消息,
		// 只有gate重启了,或者断线重连, 才会重新有此消息,
		// 所以要down后恢复, 需要持久化p.ClientShouldJoin.

		// 分配给指定的actor
		if _, ok := p.ClientShouldJoin[msg.Uid]; !ok {
			shouldJoinRoomId := "1"
			p.ClientShouldJoin[msg.Uid] = shouldJoinRoomId
		}
		
		pid, err := p.CreateOrFindGameActor(p.ClientShouldJoin[msg.Uid])
		if err != nil {
			log.ErrorT(TAG, err)
			return
		}
		pid.Request(msg, ctx.Sender())
	case *pbgo.ClientMessageReq:
		// 每一个由客户端发送的消息

		// 分配给指定的actor
		if roomId, ok := p.ClientShouldJoin[msg.Uid]; ok {
			pid, err := p.CreateOrFindGameActor(roomId)
			if err != nil {
				log.ErrorT(TAG, err)
				return
			}

			pid.Request(msg, ctx.Sender())
		} else {
			log.ErrorT(TAG, "Uid", msg.Uid, "not conned")
		}
	case *pbgo.ClientDisconnectReq:
		// 当客户端关闭连接的消息

		if roomId, ok := p.ClientShouldJoin[msg.Uid]; ok {
			pid, err := p.CreateOrFindGameActor(roomId)
			if err != nil {
				log.ErrorT(TAG, err)
				return
			}
			pid.Request(msg, ctx.Sender())

			delete(p.ClientShouldJoin, msg.Uid)
		}
	case *actor.Started:
		p.Recover()
	}
}

// 恢复Recover
func (p *GameActorManager) Recover() {
	p.ClientShouldJoin = map[string]string{}

	// debug
	shouldJoinRoomId := "1"
	p.ClientShouldJoin["1"] = shouldJoinRoomId
}

func (p *GameActorManager) CreateOrFindGameActor(roomId string) (pid *actor.PID, err error) {
	if pid, ok := p.Actors[roomId]; ok {
		return pid, nil
	}

	pid, err = actor.SpawnNamed(actor.FromInstance(NewGameActor()), roomId)
	if err == nil {
		p.Actors[roomId] = pid
	}
	return
}

func NewGameActorManager() *GameActorManager {
	return &GameActorManager{
		ClientShouldJoin: map[string]string{},
		Actors:           map[string]*actor.PID{},
	}
}

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
		// 分配给指定的actor
		if _, ok := p.ClientShouldJoin[msg.Uid]; !ok {
			shouldJoinRoomId := "1"
			p.ClientShouldJoin[msg.Uid] = shouldJoinRoomId
		}
		log.InfoT(TAG, "user", msg.Uid, "conned")

	case *pbgo.ClientMessageReq:
		// 分配给指定的actor
		if roomId, ok := p.ClientShouldJoin[msg.Uid]; ok {
			pid, err := p.CreateOrFindGameActor(roomId)
			if err != nil {
				log.ErrorT(TAG, err)
				return
			}

			pid.Tell(msg)
		} else {
			log.ErrorT(TAG, "Uid", msg.Uid, "not conned")
		}
	case *pbgo.ClientCloseRsq:
		if roomId, ok := p.ClientShouldJoin[msg.Uid]; ok {
			pid, err := p.CreateOrFindGameActor(roomId)
			if err != nil {
				log.ErrorT(TAG, err)
				return
			}
			pid.Tell(msg)

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

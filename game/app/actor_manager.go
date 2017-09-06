package app

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bysir-zl/game-frame-probe/common/pbgo"
	"github.com/bysir-zl/bygo/log"
	"github.com/bysir-zl/game-frame-probe/common/client_msg"
	"github.com/bysir-zl/game-frame-probe/common"
	"github.com/bysir-zl/game-frame-probe/common/util"
	"time"
)

type GameActorManager struct {
	// 客户端应该进入哪个Actor, 在通信之前发送消息让manager判断客户端应该加入哪个actor(比如加入某ID游戏局), 在之后的通信就不需要了; 为了down机恢复,需要持久化此字段
	ClientToActor *util.MapStorage      // uid=>roomId
	Actors        map[string]*actor.PID // 所有已经生成的actor, name=>pid
	GatePid       *actor.PID
}

func (p *GameActorManager) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *pbgo.GateConnectReq:
		// gate上线
		p.GatePid = ctx.Sender()

	case *pbgo.ClientConnectReq:
		// 在第一次客户端向这个节点发消息的时候会连接这个节点, 
		// 如果本节点重启后, 也不再有此消息,
		// 只有gate重启了,或者断线重连, 才会重新有此消息,

		// 如果要down后恢复, 需要持久化p.ClientToActor.
	case *pbgo.ClientMessageReq:
		// 每一个由客户端发送的消息

		// 分配给指定的actor
		if actorId, ok := p.ClientToActor.Get(msg.Uid); ok {
			pid, isNew, err := p.CreateOrFindGameActor(ctx.Self(), actorId)
			if err != nil {
				log.ErrorT(TAG, err)
				return
			}
			// 如果是新生产的actor(重启后会发生这个情况)
			// 重新发送一次连接消息
			if isNew {
				pid.Request(&pbgo.ClientConnectReq{Uid: msg.Uid}, ctx.Sender())
			}
			pid.Request(msg, ctx.Sender())
		} else {
			// 如果没有为client分配过actor, 那么重新分配

			// 根据客户端消息加入GameActor
			pt := client_msg.GetProto(msg.Body)
			actorId := pt.Body
			switch pt.Cmd {
			case common.CMD_JoinRoom:
				p.ClientToActor.Set(msg.Uid, actorId)

				pid, _, err := p.CreateOrFindGameActor(ctx.Self(), actorId)
				if err != nil {
					log.ErrorT(TAG, err)
					return
				}

				// 发送一次连接消息
				pid.Request(&pbgo.ClientConnectReq{Uid: msg.Uid}, ctx.Sender())
				pid.Request(msg, ctx.Sender())
			default:
				log.ErrorT(TAG, "Uid", msg.Uid, "need conn, but cmd is", pt.Cmd, ", want CMD_JoinRoom")
			}
		}
	case *pbgo.ClientDisconnectReq:
		// 当客户端关闭连接的消息

		if actorId, ok := p.ClientToActor.Get(msg.Uid); ok {
			pid, _, err := p.CreateOrFindGameActor(ctx.Self(), actorId)
			if err != nil {
				log.ErrorT(TAG, err)
				return
			}
			pid.Request(msg, ctx.Sender())

			p.ClientToActor.Del(msg.Uid)
		}
	case *ActorStopMsg:
		delete(p.Actors, msg.ActorId)
	case *actor.Started:
		p.Recover()

		go func() {
			for range time.Tick(time.Second * 6) {
				
			}
			// 检查网关状态
			// 网关死了以后, 网关保存的客户端PID就失效了, 任何Actor的数据发送都会失败, 应当关闭所有actor
			_, err := p.GatePid.RequestFuture(&pbgo.GatePing{}, time.Second*3).Result()
			if err != nil {
				return
			}
		}()
	}
}

// 恢复Recover
func (p *GameActorManager) Recover() {
	p.ClientToActor.Load()
}

func (p *GameActorManager) CreateOrFindGameActor(manager *actor.PID, actorId string) (pid *actor.PID, isCreate bool, err error) {
	if pid, ok := p.Actors[actorId]; ok {
		return pid, false, nil
	}

	pid, err = actor.SpawnNamed(actor.FromInstance(NewGameActor(manager, actorId)), actorId)
	if err == nil {
		p.Actors[actorId] = pid
		isCreate = true
	}
	return
}

func NewGameActorManager() *GameActorManager {
	return &GameActorManager{
		ClientToActor: util.NewMapStorage(id),
		Actors:        map[string]*actor.PID{},
	}
}

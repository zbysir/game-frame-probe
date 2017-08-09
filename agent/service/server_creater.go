package service

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"errors"
	"github.com/bysir-zl/game-frame-probe/proto/pbgo"
	"time"
)

// unused
func GetServiceActor(uid string, types string) (pid *actor.PID, err error) {
	serverManager, ok := stdServerGroups.SelectServer(types)
	if !ok {
		err = errors.New("404")
		return
	}

	rsp, err := serverManager.RequestFuture(&pbgo.GetServerClientActorReq{
		Uid: uid,
	}, time.Second*3).Result()
	if err != nil {
		return
	}

	pid = (rsp.(*pbgo.GetServerClientActorRsp)).Pid
	return
}

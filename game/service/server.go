package service

import (
	"time"
	"github.com/AsynkronIT/protoactor-go/actor"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/bysir-zl/game-frame-probe/common/service"
	"github.com/bysir-zl/bygo/util/discovery"
)
const TAG = "game"

var (
	id   string = "game/1"
	addr string = "127.0.0.1"
	port int    = 8090
)

func Run() {
	_, err := serverNode()
	if err != nil {
		panic(err)
	}

	// 注册服务
	lease, err := service.Discoverer.RegisterService(&discovery.Server{
		Id:      id,
		Name:    id,
		Address: addr,
		Port:    port,
	})
	if err != nil {
		panic(err)
	}

	for range time.Tick(time.Second * 5) {
		service.Discoverer.UpdateServerTTL(lease)
	}
}

func serverNode() (pid *actor.PID, err error) {
	nodeAddr := fmt.Sprintf("%s:%d", addr, port)

	remote.Start(nodeAddr)
	props := actor.FromInstance(NewGameActor())
	pid, err = actor.SpawnNamed(props, id)
	return
}

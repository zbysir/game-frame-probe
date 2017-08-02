package common

import (
	"testing"
	"github.com/bysir-zl/bygo/log"
	"time"
	"github.com/bysir-zl/game-frame-probe/common/service"
)

var consulAddr = "127.0.0.1:8500"

func TestGetServer(t *testing.T) {
	e := service.NewManagerEtcd()

	tick := time.Tick(1 * time.Second)
	for {
		select {
		case <-tick:
			ser, err := e.GetServers()
			if err != nil {
				t.Fatal(err)
			}
			for _, s := range ser {
				log.InfoT("test", s)
			}
		}
	}
}

func TestUpdateTTl(t *testing.T) {
	e := service.NewManagerEtcd()
	err := e.UpdateServerTTL("7587823890455770676")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRegisterService(t *testing.T) {
	e := service.NewManagerEtcd()
	leaseId,err:= e.RegisterService(&service.Server{
		Id:      "game-1",
		Name:    "game-1",
		Port:    8090,
		Address: "127.0.0.1",
	})
	if err != nil {
		t.Fatal(err)
	}
	
	time.Sleep(10*time.Second)
	err = e.UpdateServerTTL(leaseId)
	if err != nil {
		t.Fatal(err)
	}

}

func TestUnRegisterService(t *testing.T) {
	e := service.NewManagerEtcd()
	err := e.UnRegisterService("game-1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestWatch(t *testing.T) {
	e := service.NewManagerEtcd()
	e.WatchServer(func(server *service.Server, change service.ServerChange) {
		log.InfoT("test", server, change)
	})
	select {}
}

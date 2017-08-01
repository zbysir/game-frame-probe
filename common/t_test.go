package common

import (
	"testing"
	"github.com/bysir-zl/game-frame-probe/common/app"
	"github.com/bysir-zl/bygo/log"
	"time"
)

var consulAddr = "127.0.0.1:8500"

func TestGetServer(t *testing.T) {
	app.Init()

	tick := time.Tick(1 * time.Second)
	for {
		select {
		case <-tick:
			ser, err := app.GetServices()
			if err != nil {
				t.Fatal(err)
			}
			for _, s := range ser {
				log.InfoT("test",s)
			}
		}
	}
}

func TestUpdateTTl(t *testing.T) {
	err := app.UpdateServerTTL("game-1", "pass")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRegisterService(t *testing.T) {
	err := app.RegisterService(&app.Service{
		Id:      "game-1",
		Name:    "game-1",
		Address: "127.0.0.1",
		Port:    8000,
	}, "30s")
	if err != nil {
		t.Fatal(err)
	}
}

func TestUnRegisterService(t *testing.T) {
	err := app.UnRegisterService("game-3")
	if err != nil {
		t.Fatal(err)
	}
}

func TestWatch(t *testing.T) {
	change, err := app.Watch("checks", "services")
	if err != nil {
		t.Fatal(err)
	}
	for {
		c := <-change
		log.InfoT("test", c)
	}
}

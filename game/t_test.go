package main

import (
	"testing"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/bysir-zl/game-frame-probe/common/pbgo"
)

func TestSendMessage(t *testing.T) {
	remote.Start("127.0.0.1:0")

	server := actor.NewPID("127.0.0.1:8090", "game/1")
	//server.Tell(&pbgo.ClientConnectReq{Uid: "1"})
	server.Tell(&pbgo.ClientMessageReq{Uid: "1", Body: []byte("hello")})
	server.Tell(&pbgo.ClientMessageReq{Uid: "2", Body: []byte("hello")})
	server.Tell(&pbgo.ClientMessageReq{Uid: "1", Body: []byte("hello agent")})

	select {}
}

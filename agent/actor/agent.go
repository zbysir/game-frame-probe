package actor

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bysir-zl/game-frame-probe/proto/pbgo"
	"fmt"
	"reflect"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/bysir-zl/bygo/log"
	"github.com/bysir-zl/hubs/core/hubs"
	"github.com/bysir-zl/hubs/core/net/listener"
	"github.com/bysir-zl/hubs/core/net/conn_wrap"
	"encoding/json"
)

type AgentActor struct {
}

func (p *AgentActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *pbgo.AgentConnect:
		// 当有节点连接上这个网关
		fmt.Println(msg.Sender)
		props := actor.FromInstance(&AgentActor{})
		pid := actor.Spawn(props)
		msg.Sender.Tell(&pbgo.AgentConnected{Message: "Welcome!", Server: pid})
	case *pbgo.AgentForwardToSvr:
		fmt.Println(msg.ServerName)

	default:
		log.Info(reflect.TypeOf(msg).String())
	}
}

func Server() {
	pid, err := serverNode()
	if err != nil {
		panic(err)
	}

	serverCli(pid)
}

func serverNode() (pid *actor.PID, err error) {
	nodeAddr := "127.0.0.1:8080"

	remote.Start(nodeAddr)
	props := actor.FromInstance(&AgentActor{})
	pid, err = actor.SpawnNamed(props, "agent")
	return
}

func clientHandle(server *hubs.Server, conn conn_wrap.Interface) {
	for {
		bs, err := conn.Read()
		if err != nil {
			return
		}
		type Proto struct {
			Cmd  int
			Body string
		}
		p := Proto{}
		json.Unmarshal(bs, &p)
		switch p.Cmd {
		case 1:
			msg := pbgo.AgentForwardToSvr{ServerName: "game", Body: []byte(p.Body), Uid: "1"}
			dPid.Tell(&msg)
		}
	}
}

var dPid *actor.PID

func serverCli(pid *actor.PID) {
	dPid = pid
	addr := "127.0.0.1:8081"

	hubServer := hubs.New(addr, listener.NewWs(), clientHandle)
	hubServer.Run()
}

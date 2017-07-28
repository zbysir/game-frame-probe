package actor

import (
	"testing"
	"github.com/AsynkronIT/protoactor-go/actor"
	"log"
	"github.com/AsynkronIT/protoactor-go/remote"
	reflect "reflect"
	"github.com/bysir-zl/game-frame-probe/proto/pbgo"
)

func TestRun(t *testing.T) {
	Run()
}

func TestTell(t *testing.T) {
	remote.Start("127.0.0.1:0")

	server := actor.NewPID("127.0.0.1:8080", "agent")

	props := actor.FromFunc(func(context actor.Context) {
		switch msg := context.Message().(type) {
		case *pbgo.Connect:
			log.Println(msg.Sender)
		case *pbgo.SayResponse:
			log.Printf("%v: %v", msg.UserName, msg.Message)
		case *pbgo.NickResponse:
			log.Printf("%v is now known as %v", msg.OldUserName, msg.NewUserName)
		default:
			log.Print(reflect.TypeOf(msg))
		}
	})
	client := actor.Spawn(props)

	server.Tell(&pbgo.Connect{
		Sender: client,
	})

	<-(chan int)(nil)
}

func TestTell2(t *testing.T) {
	remote.Start("127.0.0.1:0")

	server := actor.NewPID("127.0.0.1:8080", "agent")

	props := actor.FromFunc(func(context actor.Context) {
		switch msg := context.Message().(type) {
		case *pbgo.Connected:
			log.Println(msg.Message)
		case *pbgo.SayResponse:
			log.Printf("%v: %v", msg.UserName, msg.Message)
		case *pbgo.NickResponse:
			log.Printf("%v is now known as %v", msg.OldUserName, msg.NewUserName)
		default:
			log.Print(reflect.TypeOf(msg))
		}
	})
	client := actor.Spawn(props)

	server.Tell(&pbgo.Connect{
		Sender: client,
	})

	<-(chan int)(nil)
}

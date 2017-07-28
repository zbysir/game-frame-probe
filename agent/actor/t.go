package actor

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"k5star/chesscards/util/log"
	"reflect"
	"github.com/bysir-zl/game-frame-probe/proto/pbgo"
)

type Response struct {
	count int
}

type MyActor struct {
	count int
}

func (p *MyActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *pbgo.Connect:
		fmt.Println(msg.Sender)
		msg.Sender.Tell(&pbgo.Connected{Message: "Welcome!"})
		
	default: 
		log.Print(reflect.TypeOf(msg))
	}
}

func Run() {
	remote.Start("127.0.0.1:8080")

	props := actor.FromInstance(&MyActor{})
	actor.SpawnNamed(props, "agent")
	
	//time.Sleep(1 * time.Second)
	<-(chan int)(nil)
}

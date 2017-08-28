package app

import "github.com/AsynkronIT/protoactor-go/actor"


type Player struct {
	Uid string `json:"id"`
	*actor.PID `json:"-"`
}

type Broad struct {
	Uid string `json:"id"`
	Body string `json:"body"`
}



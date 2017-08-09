package service

import "github.com/AsynkronIT/protoactor-go/actor"

type Point struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type Player struct {
	Uid string `json:"id"`
	Angle float32 `json:"angle"`
	Speed float32 `json:"speed"`
	
	*actor.PID `json:"-"`

	Point 
}



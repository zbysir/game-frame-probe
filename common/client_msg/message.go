package client_msg

import "encoding/json"

type Proto struct {
	Cmd  int `json:"cmd"`
	Body string `json:"body"`
}

type MsgMove struct {
	Angle float32 `json:"angle"`
	Speed float32 `json:"speed"`
}
type MsgJoinRoom struct {
	RoomId int64 `json:"room_id"`
}

func GetProto(data []byte) *Proto {
	p := Proto{}
	json.Unmarshal(data, &p)
	return &p
}

func GetBody(data []byte) []byte {
	p := Proto{}
	json.Unmarshal(data, &p)
	return []byte(p.Body)
}

// {cmd:0,body:'{angle,speed}'}
func GetMove(data []byte) *MsgMove {
	m := MsgMove{}
	json.Unmarshal(GetBody(data), &m)
	return &m
}

// {cmd:0,body:'{room_id}'}
func GetJoinRoom(data []byte) *MsgJoinRoom {
	m := MsgJoinRoom{}
	json.Unmarshal(GetBody(data), &m)
	return &m
}

func NewProto(cmd int,body []byte)[]byte{
	bs,_:=json.Marshal(&Proto{Cmd:cmd,Body:string(body)})
	return bs
}
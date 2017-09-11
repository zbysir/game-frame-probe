package client_msg

import "encoding/json"

type Proto struct {
	Cmd  int `json:"cmd"`
	Body string `json:"body"`
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

func NewProto(cmd int, body []byte) []byte {
	bs, _ := json.Marshal(&Proto{Cmd: cmd, Body: string(body)})
	return bs
}

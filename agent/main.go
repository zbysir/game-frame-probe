package main

import (
	"github.com/bysir-zl/game-frame-probe/agent/act"
	"github.com/bysir-zl/game-frame-probe/agent/service"
)

// 网关
// 分发消息

// length[4,uint32] cmd[4] bytes[~]

func main() {
	service.Run()
}

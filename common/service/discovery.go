package service

import "github.com/bysir-zl/bygo/util/discovery"

var Discoverer = discovery.NewEtcd([]string{"localhost:2379"}, 30)

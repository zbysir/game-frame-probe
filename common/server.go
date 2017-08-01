package common

import (
	"github.com/hashicorp/consul/api"
	"github.com/bysir-zl/bygo/log"
	"github.com/hashicorp/consul/watch"
)

func GetServerNode(addr string) (err error) {
	consulConf := api.DefaultConfig()
	consulConf.Address = addr
	client, err := api.NewClient(consulConf)
	if err != nil {
		return
	}
	services, _, err := client.Catalog().Services(&api.QueryOptions{})
	if err != nil {
		return
	}
	
	for name := range services {
		servicesData, _, err := client.Health().Service(name, "", false,
			&api.QueryOptions{})
		if err != nil {
			log.ErrorT("e", err)
			continue
		}

		for _, entry := range servicesData {
			log.InfoT("consul", entry.Service,entry.Node,entry.Checks.AggregatedStatus())
		}

	}
	return
}

func RegisterService(addr string) (err error) {
	consulConf := api.DefaultConfig()
	consulConf.Address = addr
	client, err := api.NewClient(consulConf)
	if err != nil {
		return
	}

	server := api.AgentServiceRegistration{
		ID:      "game-1",
		Address: "127.0.0.1",
		Name:    "game-1",
		Port:    8090,
		Tags:    []string{"game-1"},
		Check: &api.AgentServiceCheck{
			TTL: "30s",
		},
	}

	agent := client.Agent()
	err = agent.ServiceRegister(&server)
	if err != nil {
		return
	}
	err = agent.PassTTL("service:game-1","")
	if err != nil {
		return 
	}
	
	return
}

func Watch(addr string) {
	params := map[string]interface{}{
		"type": "checks",
	}
	plan, err := watch.Parse(params)
	if err != nil {
		return
	}
	plan.Handler = func(id uint64, raw interface{}) {
		checks:=raw.([]*api.HealthCheck)
		log.InfoT("watch", id, raw,checks[0])
	}

	plan.Run(addr)
}

package app

import (
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/watch"
	"github.com/bysir-zl/bygo/log"
	"reflect"
	"sync"
	"errors"
)

// 用来做服务发现

var stdCli *api.Client

type Service struct {
	Id      string
	Name    string
	Address string
	Port    int
	Status  string
}

var stdServices map[string]*Service
var lock sync.Mutex
var onChangedFunc func(services map[string]*Service)

func GetServices() (services map[string]*Service, err error) {
	if stdServices == nil {
		err = errors.New("not init")
		return
	}
	
	return stdServices, nil
}

func updateServicesById(name []string, passOnly bool) {
	lock.Lock()
	defer lock.Unlock()

	stdServices = map[string]*Service{}
	for _, name := range name {
		servicesData, _, e := stdCli.Health().Service(name, "", passOnly,
			&api.QueryOptions{})
		if e != nil {
			return
		}

		for _, entry := range servicesData {
			stdServices[entry.Service.ID] = &Service{
				Address: entry.Service.Address,
				Port:    entry.Service.Port,
				Name:    entry.Service.Service,
				Id:      entry.Service.ID,
				Status:  entry.Checks.AggregatedStatus(),
			}
		}
	}

	if onChangedFunc !=nil{
		onChangedFunc(stdServices)
	}

	return
}

func UpdateServerTTL(serviceId string, status string) (err error) {
	err = stdCli.Agent().UpdateTTL("service:"+serviceId, "", status)
	return
}
func OnChangedFunc(onChanged func(services map[string]*Service)){
	onChangedFunc = onChanged
}

func RegisterService(service *Service, ttl string) (err error) {
	server := &api.AgentServiceRegistration{
		ID:      service.Id,
		Address: service.Address,
		Name:    service.Name,
		Port:    service.Port,
		Tags:    []string{},
		Check: &api.AgentServiceCheck{
			TTL: ttl,
		},
	}
	//UpdateServerTTL(server.ID,"pass")
	err = stdCli.Agent().ServiceRegister(server)
	return
}

func UnRegisterService(serviceId string) (err error) {
	err = stdCli.Agent().ServiceDeregister(serviceId)
	return
}

type Change struct {
	Type string
	Obj  interface{}
}

func Watch(types ...string) (changeC chan *Change, err error) {
	changeC = make(chan *Change, 100)
	for _, t := range types {
		go func(t string) {
			plan, e := watch.Parse(map[string]interface{}{
				"type": t,
			})
			if e != nil {
				return
			}
			plan.Handler = func(id uint64, raw interface{}) {
				switch data := raw.(type) {
				case map[string][]string:
					names := []string{}
					for name := range data {
						names = append(names, name)
					}
					updateServicesById(names, false)

				case []*api.HealthCheck:
					for _, c := range data {
						log.InfoT("ttt", c)
						if ser, ok := stdServices[c.ServiceID]; ok {
							ser.Status = c.Status
						}
					}

				default:
					log.InfoT("eee", reflect.TypeOf(raw))
				}

				changeC <- &Change{
					Type: plan.Type,
					Obj:  raw,
				}
			}
			plan.Run(consulAddr)
		}(t)
	}
	return
}

func watchService(types ...string) {
	isReady := false
	readyChan := make(chan int)
	for _, t := range types {
		go func(t string) {
			plan, e := watch.Parse(map[string]interface{}{
				"type": t,
			})
			if e != nil {
				return
			}
			plan.Handler = func(id uint64, raw interface{}) {
				switch data := raw.(type) {
				case map[string][]string:
					names := []string{}
					for name := range data {
						names = append(names, name)
					}
					updateServicesById(names, false)
				case []*api.HealthCheck:
					names := []string{}
					for _, name := range data {
						names = append(names, name.ServiceID)
					}
					updateServicesById(names, false)
				}
				if !isReady {
					isReady = true
					close(readyChan)
				}
			}
			plan.Run(consulAddr)
		}(t)
	}
	select {
	case <-readyChan:
	}
	return
}

func Init() {
	watchService("services", "checks")
}

var consulAddr = "127.0.0.1:8500"

func init() {
	consulConf := api.DefaultConfig()
	consulConf.Address = consulAddr
	cli, err := api.NewClient(consulConf)
	if err != nil {
		panic(err)
	}
	stdCli = cli
}

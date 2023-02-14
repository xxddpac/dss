package consul

import (
	"dss/core/host"
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
)

var (
	apiClient   *api.Client
	ServiceName string
	Port        int
)

type Config struct {
	Host       string
	Port       int
	DataCenter string
}

type ServiceRegistration struct {
	ID      string                 `json:"ID"`
	Name    string                 `json:"Name"`
	Tags    []string               `json:"Tags,omitempty"`
	Address string                 `json:"Address"`
	Port    int                    `json:"Port"`
	Check   *api.AgentServiceCheck `json:"Check,omitempty"`
	Meta    map[string]string      `json:",omitempty"`
}

func newClient(cfg *Config) (*api.Client, error) {
	_config := api.DefaultConfig()
	_config.Datacenter = cfg.DataCenter
	_config.Address = fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	return api.NewClient(_config)
}

func Init(cfg *Config) error {
	cli, err := newClient(cfg)
	if err != nil {
		return fmt.Errorf("new consul client err: %s", err.Error())
	}
	apiClient = cli
	return nil
}

func Register() {
	var (
		serviceID string
		localIP   = host.LocalIP()
	)
	serviceID = fmt.Sprintf("%s-%s", ServiceName, localIP)
	svr := ServiceRegistration{
		ID:      serviceID,
		Name:    ServiceName,
		Tags:    []string{ServiceName},
		Address: localIP,
		Port:    Port,
	}
	svr.Check = &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d%s", svr.Address, svr.Port, "/ping"),
		Timeout:                        "3s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "30s",
	}
	if err := serviceRegister(&svr); err != nil {
		log.Fatalf("failed register service to consul, err:%v", err)
	}
}

func serviceRegister(svr *ServiceRegistration) error {
	registration := new(api.AgentServiceRegistration)
	registration.ID = svr.ID
	registration.Name = svr.Name
	registration.Port = svr.Port
	registration.Tags = svr.Tags
	registration.Address = svr.Address
	registration.Check = svr.Check
	registration.Meta = svr.Meta
	if err := apiClient.Agent().ServiceRegister(registration); nil != err {
		return fmt.Errorf("register server err: %s", err.Error())
	}
	return nil
}

func Deregister(serviceID string) error {
	return apiClient.Agent().ServiceDeregister(serviceID)
}

func Health() *api.Health {
	return apiClient.Health()
}

func Close() {
	apiClient = nil
}

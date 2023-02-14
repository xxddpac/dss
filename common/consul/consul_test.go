package consul

import (
	"dss/core/host"
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"testing"
)

var (
	fakeConsulConfig = &Config{
		Port:       8500,
		Host:       "10.90.81.179",
		DataCenter: "DC1",
	}
)

func TestConsul(t *testing.T) {
	host.RefreshHost()
	if err := Init(fakeConsulConfig); err != nil {
		t.Fatal("Init consul failed", zap.Error(err))
	}
	svr := ServiceRegistration{
		ID:      "test_dss",
		Name:    "test_dss",
		Tags:    []string{"test_dss"},
		Meta:    map[string]string{"weight": "100"},
		Address: host.LocalIP(),
		Port:    9091,
	}
	svr.Check = &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d%s", svr.Address, svr.Port, "/ping"),
		Timeout:                        "3s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "30s",
	}
	if err := serviceRegister(&svr); nil != err {
		t.Fatal(err)
	}
}

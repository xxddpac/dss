package discover

import (
	"context"
	"dss/common/consul"
	"dss/common/log"
	"fmt"
	"github.com/hashicorp/consul/api"
	"net"
	"runtime/debug"
	"strconv"
	"time"
)

const (
	WaitTime   = 60 * time.Second
	MaxRetries = 50
)

const (
	Ok = iota + 1
	Ng
)

type Next func() string

type ServiceDiscovery struct {
	serviceName string
	lastIndex   uint64
	curIndex    int
	activeList  []string
	sig         chan int
	retries     int
}

func NewServiceDiscovery(ctx context.Context, serviceName string) *ServiceDiscovery {
	var s = &ServiceDiscovery{
		serviceName: serviceName,
		activeList:  []string{},
		sig:         make(chan int),
	}
	go s.watch(ctx)
	return s
}

func (s *ServiceDiscovery) watch(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			log.Errorf("panic err:%v", fmt.Sprint(err))
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			items, meta, err := consul.Health().Service(s.serviceName, "", true, &api.QueryOptions{
				WaitIndex: s.lastIndex,
				WaitTime:  WaitTime,
			})
			if s.retries > MaxRetries {
				log.Errorf("discover will shutdown, err: failed to obtain the available address for %d consecutive times", MaxRetries)
				s.sig <- Ng
				return
			}
			if err != nil || len(items) == 0 {
				s.retries++
				log.Errorf("failed discover service %s, err:%v", s.serviceName, err)
				time.Sleep(WaitTime)
				continue
			}
			s.update(items)
			s.lastIndex = meta.LastIndex
		}
	}
}

func (s *ServiceDiscovery) Wait() <-chan int {
	return s.sig
}

func (s *ServiceDiscovery) update(conf []*api.ServiceEntry) {
	s.activeList = []string{}
	for _, entry := range conf {
		addr := net.JoinHostPort(entry.Service.Address, strconv.Itoa(entry.Service.Port))
		s.activeList = append(s.activeList, addr)
	}
	s.sig <- Ok
}

func (s *ServiceDiscovery) Next() string {
	if len(s.activeList) == 0 {
		return ""
	}
	lens := len(s.activeList)
	if s.curIndex >= lens {
		s.curIndex = 0
	}
	curAddr := s.activeList[s.curIndex]
	s.curIndex = (s.curIndex + 1) % lens
	return curAddr
}

func (s *ServiceDiscovery) Get() string {
	return s.Next()
}

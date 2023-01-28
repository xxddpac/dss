package dao

import (
	"fmt"
	"dss/common/redis"
	"testing"
	"time"
)

var (
	fakeRedisConfig = &redis.Config{
		Network: "tcp",
		Addr:    "10.101.191.106:6379",
		Passwd:  "",
		DB:      13,
	}
	key = "test_port_scan"
)

func TestRedis(t *testing.T) {
	if err := redis.Init(fakeRedisConfig); err != nil {
		t.Fatal(err)
	}
	go func() {
		for i := 1; i <= 5; i++ {
			if err := Redis.LPush(key, fmt.Sprintf("test_%v", i)); err != nil {
				t.Error(err)
				return
			}
		}
	}()
	go func() {
		for {
			if msg, err := Redis.BRPop(key, 0*time.Second); err == nil {
				fmt.Printf("Received msg from node_1 --->%v \n", msg)
			}
		}
	}()
	go func() {
		for {
			if msg, err := Redis.BRPop(key, 0*time.Second); err == nil {
				fmt.Printf("Received  msg from node_2 --->%v \n", msg)
			}
		}
	}()
	time.Sleep(2 * time.Second)
}

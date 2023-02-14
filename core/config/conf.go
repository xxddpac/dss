package config

import (
	"dss/common/consul"
	"dss/common/log"
	"dss/common/mongo"
	"dss/common/redis"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"time"
)

var (
	CoreConf *config
)

func Init(conf string) {
	_, err := toml.DecodeFile(conf, &CoreConf)
	if err != nil {
		fmt.Printf("Err %v", err)
		os.Exit(1)
	}
}

type config struct {
	Log         log.Config
	Redis       redis.Config
	Mongo       mongo.Config
	Consumer    Consumer
	Producer    Producer
	Consul      consul.Config
	Mode        string
	ServiceName string
	MaxWorkers  int
	MaxQueue    int
	GrpcPort    int
}

type Consumer struct {
	Port    int
	TimeOut time.Duration
}

type Producer struct {
	Port              int
	WorkChatUploadUrl string
	WorkChatBotUrl    string
}

package config

import (
	"dss/common/consul"
	"dss/common/log"
	"dss/common/mongo"
	"dss/common/redis"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
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
	Port  int
	Pprof *pprof
}

type Producer struct {
	Port              int
	WorkChatUploadUrl string
	WorkChatBotUrl    string
	Pprof             *pprof
}

type pprof struct {
	Enable bool
	Port   int
}

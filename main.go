package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"goportscan/common/log"
	"goportscan/common/mongo"
	"goportscan/common/redis"
	"goportscan/core/config"
	"goportscan/core/router"
	"goportscan/core/scan"
	"goportscan/core/server"
	"os"
	"runtime"
)

var (
	cfg string
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	flag.StringVar(&cfg, "conf", "", "server config [toml]")
	flag.Parse()
	if len(cfg) == 0 {
		fmt.Println("config is empty")
		os.Exit(0)
	}
	//init config
	config.Init(cfg)
	conf := config.CoreConf
	//init log
	log.Init(&conf.Log)
	//init redis
	if err := redis.Init(&conf.Redis); err != nil {
		log.Fatal("Init redis failed", zap.Error(err))
	}
	//init mongo
	if err := mongo.Init(&conf.Mongo); err != nil {
		log.Fatal("Init mongo failed", zap.Error(err))
	}
	//init port scan
	scan.Init(conf.Server.MaxWorkers, conf.Server.MaxQueue, log.Logger())
	//init gin web server
	gin.SetMode(conf.Server.Mode)
	if err := server.Run(router.NewHttpRouter()); err != nil {
		log.Error("server run error", zap.Error(err))
	}
}

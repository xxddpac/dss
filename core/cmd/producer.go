package cmd

import (
	"dss/common/consul"
	"dss/common/log"
	"dss/common/mongo"
	"dss/common/redis"
	"dss/core/config"
	"dss/core/global"
	"dss/core/grpc/producer"
	"dss/core/host"
	"dss/core/management"
	"dss/core/pprof"
	"dss/core/router"
	"dss/core/server"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	_ "net/http/pprof"
	"os"
)

func Producer() *cobra.Command {
	var (
		cfg string
	)
	cmdProducer := &cobra.Command{
		Use:   "producer",
		Short: "Start Run Security Scan Producer",
		Run: func(cmd *cobra.Command, args []string) {
			if len(cfg) == 0 {
				_ = cmd.Help()
				os.Exit(0)
			}
			config.Init(cfg)
			conf := config.CoreConf
			consul.ServiceName = conf.ServiceName
			consul.Port = conf.Producer.Port
			log.Init(&conf.Log)
			host.RefreshHost()
			producer.Grpc()
			go pprof.Pprof(conf.Producer.Pprof.Enable, conf.Producer.Pprof.Port)
			if err := redis.Init(&conf.Redis); err != nil {
				log.Fatal("Init redis failed", zap.Error(err))
			}
			if err := mongo.Init(&conf.Mongo); err != nil {
				log.Fatal("Init mongo failed", zap.Error(err))
			}
			if err := consul.Init(&conf.Consul); err != nil {
				log.Fatal("Init consul failed", zap.Error(err))
			}
			go management.RunTimeTaskStatusCheck()
			consul.Register()
			gin.SetMode(conf.Mode)
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			if err := server.Run(router.NewHttpRouter(), global.Producer); nil != err {
				log.Error("server run error", zap.Error(err))
			}
		},
	}
	cmdProducer.Flags().StringVarP(&cfg, "conf", "c", "", "server config [toml]")
	return cmdProducer
}

package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"goportscan/common/log"
	"goportscan/common/mongo"
	"goportscan/common/redis"
	"goportscan/core/config"
	"goportscan/core/global"
	"goportscan/core/router"
	"goportscan/core/server"
	"os"
)

func Producer() *cobra.Command {
	var (
		cfg string
	)
	cmdConsumer := &cobra.Command{
		Use:   "producer",
		Short: "Start Run Port Scan Producer",
		Run: func(cmd *cobra.Command, args []string) {
			if len(cfg) == 0 {
				_ = cmd.Help()
				os.Exit(0)
			}
			config.Init(cfg)
			conf := config.CoreConf
			log.Init(&conf.Log)
			if err := redis.Init(&conf.Redis); err != nil {
				log.Fatal("Init redis failed", zap.Error(err))
			}
			if err := mongo.Init(&conf.Mongo); err != nil {
				log.Fatal("Init mongo failed", zap.Error(err))
			}
			gin.SetMode(conf.Consumer.Mode)
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			if err := server.Run(router.NewHttpRouter(), global.Producer); nil != err {
				log.Error("server run error", zap.Error(err))
			}
		},
	}
	cmdConsumer.Flags().StringVarP(&cfg, "conf", "c", "", "server config [toml]")
	return cmdConsumer
}

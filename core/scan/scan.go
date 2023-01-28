package scan

import (
	"dss/common/async"
	"dss/common/log"
	"dss/common/utils"
	"dss/core/config"
	"dss/core/dao"
	"dss/core/global"
	"dss/core/models"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net"
	"time"
)

var (
	pool    *async.WorkerPool
	timeout time.Duration
	result  = make([]scanInfo, 0)
	queue   = make(chan scanInfo, 100)
	repo    = dao.Repository{Collection: "port_scan"}
)

type scanInfo struct {
	Host     string
	Port     string
	Location string
}

func Init(maxWorkers, maxQueue int, log *zap.Logger) {
	timeout = config.CoreConf.Consumer.TimeOut
	pool = async.NewWorkerPool(maxWorkers, maxQueue, log).Run()
	go run()
	go store()
}

func Close() {
	pool.Close()
}

func run() {
	log.Info("start run port scan...")
	for {
		if msg, err := dao.Redis.BRPop(global.ScanQueue, 0*time.Second); err == nil {
			s := &scanInfo{}
			if err := json.Unmarshal([]byte(msg), s); err != nil {
				log.Error("failed unmarshal json at parse message", zap.String("msg", msg), zap.Error(err))
				continue
			}
			pool.Add(s)
		}
	}
}

func (s *scanInfo) Do() {
	client, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", s.Host, s.Port), timeout*time.Second)
	if err == nil {
		_ = client.Close()
		log.Infof("found host:%s open port:%s", s.Host, s.Port)
		queue <- *s
	}
}

func store() {
	//every minute collecting data insert into mongo
	ticker := time.NewTicker(time.Minute * 1)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			var (
				resp []interface{}
			)
			if len(result) == 0 {
				continue
			}
			for _, item := range result {
				resp = append(resp, models.ScanInsertFunc(models.Scan{
					Host:     item.Host,
					Port:     item.Port,
					Location: item.Location,
					DoneTime: time.Now().Format(utils.TimeLayout),
				}))
			}
			if err := repo.BulkWrite(resp); err != nil {
				log.Errorf("insert data to mongo err:", err)
				continue
			}
			log.Infof("success insert data to mongo,total:%d", len(resp))
			result = append(result[:0], result[len(result):]...)
		case msg := <-queue:
			result = append(result, msg)
		case <-global.Ctx.Done():
			return
		}
	}
}

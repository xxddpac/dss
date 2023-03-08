package scan

import (
	"dss/common/async"
	"dss/common/log"
	"dss/common/utils"
	"dss/core/buffer"
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
	poolForPortScan *async.WorkerPool
	//poolForDispatch *async.WorkerPool
	queue = make(chan models.Scan, 30000)
)

type scanInfo models.Scan

func Init(maxWorkers, maxQueue int, log *zap.Logger) {
	poolForPortScan = async.NewWorkerPool(maxWorkers, maxQueue, log).Run()
	//poolForDispatch = async.NewWorkerPool(maxWorkers, maxQueue, log).Run()
	go run()
}

func Close() {
	poolForPortScan.Close()
	//poolForDispatch.Close()
}

func run() {
	log.Info("start run port scan...")
	for {
		if msg, err := dao.Redis.BRPop(global.PortScanQueue, 0*time.Second); err == nil {
			s := &scanInfo{}
			if err = json.Unmarshal([]byte(msg), s); err != nil {
				log.Error("failed unmarshal json at parse message", zap.String("msg", msg), zap.Error(err))
				continue
			}
			poolForPortScan.Add(s)
		}
	}
}

func (s *scanInfo) Do() {
	item := models.Scan{
		Host:     s.Host,
		Location: s.Location,
		Port:     s.Port,
		TaskId:   s.TaskId,
		DoneTime: time.Now().Format(utils.TimeLayout),
	}
	client, err := net.DialTimeout(global.TCP, fmt.Sprintf("%v:%v", s.Host, s.Port), 2*time.Second)
	if err == nil {
		_ = client.Close()
		log.InfoF("found host:%s open port:%s", s.Host, s.Port)
		if err = dao.Repo(global.Scan).Insert(models.ScanInsertFunc(item)); err != nil {
			log.Errorf("insert scan result to mongo err:%s", err)
		}
	}
	queue <- item
}

func RunTimeUpdateTaskStatus() {
	log.InfoF("start runTimeUpdateTaskStatus...")
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			result := buffer.ReadRecords()
			tmp := make(map[string]int)
			for _, item := range result {
				tmp[item.TaskId] += 1
			}
			for taskId, count := range tmp {
				if err := dao.Redis.LuaRun(taskId, count); err != nil {
					log.Errorf("update redis err:%s", err)
				}
			}
		case task := <-queue:
			_ = buffer.WriteRecord(&task, false)
		case <-global.Ctx.Done():
			return
		}
	}
}

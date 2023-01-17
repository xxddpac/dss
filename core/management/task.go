package management

import (
	"fmt"
	"goportscan/common/log"
	"goportscan/common/utils"
	"goportscan/core/dao"
	"goportscan/core/global"
	"goportscan/core/models"
	"strings"
)

var (
	TaskManager *_TaskManager
	taskQueue   = make(chan models.Scan, 1000)
)

type _TaskManager struct {
}

func parseRule(ruleSlice []models.Rule) {
	for _, item := range ruleSlice {
		if !item.Status {
			continue
		}
		portRange := strings.Split(item.TargetPort, "-")
		portStart := utils.StrToInt(portRange[0])
		portEnd := utils.StrToInt(portRange[1])
		switch item.Type {
		case global.Single:
			for i := portStart; i <= portEnd; i++ {
				taskQueue <- models.Scan{Host: item.TargetHost, Port: fmt.Sprintf("%v", i)}
			}
		case global.Range:
			//todo
		case global.Cidr:
			ipSlice := utils.GetIpListByCidr(item.TargetHost)
			for _, ip := range ipSlice {
				for i := portStart; i <= portEnd; i++ {
					taskQueue <- models.Scan{Host: ip, Port: fmt.Sprintf("%v", i)}
				}
			}
		}
	}
}

func (*_TaskManager) Post() {
	var (
		val       string
		err       error
		ruleSlice []models.Rule
		repo      = dao.Repository{Collection: "port_scan_rule"}
	)
	if err = repo.SelectAll(&ruleSlice); err != nil {
		log.Errorf("select rule err:%v", err)
		return
	}
	go parseRule(ruleSlice)
	for {
		select {
		case task := <-taskQueue:
			val, err = utils.Marshal(task)
			if err != nil {
				log.Errorf("Marshal json to str err:%v", err)
				continue
			}
			if err = dao.Redis.LPush(global.ScanQueue, val); err != nil {
				log.Errorf("push msg to redis err:%v", err)
				continue
			}
		case <-global.Ctx.Done():
			return
		}
	}
}

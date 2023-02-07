package management

import (
	"dss/common/log"
	"dss/common/utils"
	"dss/core/dao"
	"dss/core/global"
	"dss/core/models"
	"fmt"
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
				taskQueue <- models.Scan{
					Host:     item.TargetHost,
					Port:     fmt.Sprintf("%v", i),
					Location: item.Location,
				}
			}
		case global.Range:
			//192.168.1.10-30
			ipRange := strings.Split(item.TargetHost, "-")
			ipStart := ipRange[0]
			ipSplit := strings.Split(ipStart, ".")
			ipEndLastNum := utils.StrToInt(ipRange[1])
			ipStartLastNum := utils.StrToInt(ipSplit[3])
			prefix := fmt.Sprintf("%v.%v.%v.", ipSplit[0], ipSplit[1], ipSplit[2])
			for i := ipStartLastNum; i <= ipEndLastNum; i++ {
				for p := portStart; p <= portEnd; p++ {
					taskQueue <- models.Scan{
						Host:     fmt.Sprintf("%v%v", prefix, i),
						Port:     fmt.Sprintf("%v", p),
						Location: item.Location,
					}
				}
			}
		case global.Cidr:
			ipSlice := utils.GetIpListByCidr(item.TargetHost)
			for _, ip := range ipSlice {
				for i := portStart; i <= portEnd; i++ {
					taskQueue <- models.Scan{
						Host:     ip,
						Port:     fmt.Sprintf("%v", i),
						Location: item.Location,
					}
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
	)
	if err = dao.Repo(global.ScanRule).SelectAll(&ruleSlice); err != nil {
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
			if err = dao.Redis.LPush(global.PortScanQueue, val); err != nil {
				log.Errorf("push msg to redis err:%v", err)
				continue
			}
		case <-global.Ctx.Done():
			return
		}
	}
}

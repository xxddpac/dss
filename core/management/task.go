package management

import (
	"goportscan/common/log"
	"goportscan/common/utils"
	"goportscan/core/dao"
	"goportscan/core/global"
	"goportscan/core/models"
)

var TaskManager *_TaskManager

type _TaskManager struct {
}

func parseRule() ([]models.Scan, error) {
	//todo
	return nil, nil
}

func (*_TaskManager) Post() {
	var (
		val       string
		ruleSlice []models.Rule
		err       error
		result    []models.Scan
		repo      = dao.Repository{Collection: "port_scan_rule"}
	)
	if err = repo.SelectAll(&ruleSlice); err != nil {
		log.Errorf("select rule err:%v", err)
		return
	}
	result, err = parseRule()
	if err != nil {
		log.Errorf("parse rule err:%v", err)
		return
	}
	for _, item := range result {
		val, err = utils.Marshal(item)
		if err != nil {
			log.Errorf("Marshal json to str err:%v", err)
			continue
		}
		if err = dao.Redis.LPush(global.ScanQueue, val); err != nil {
			log.Errorf("push msg to redis err:%v", err)
			continue
		}
	}

}

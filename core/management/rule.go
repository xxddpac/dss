package management

import (
	"goportscan/core/dao"
	"goportscan/core/models"
)

var (
	RuleManager *_RuleManager
	repo        = dao.Repository{Collection: "port_scan_rule"}
)

type _RuleManager struct {
}

func (*_RuleManager) Post(body models.Rule) error {
	if err := repo.Insert(models.RuleInsertFunc(body)); err != nil {
		return err
	}
	return nil
}

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

func (*_RuleManager) Delete(param models.QueryID) error {
	if err := repo.RemoveByID(dao.BsonId(param.ID)); err != nil {
		return err
	}
	return nil
}

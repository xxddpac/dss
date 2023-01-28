package management

import (
	"dss/common/utils"
	"dss/core/dao"
	"dss/core/models"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"math"
	"time"
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
	if !bson.IsObjectIdHex(param.ID) {
		return fmt.Errorf("invalid ObjectIdHex")
	}
	if err := repo.RemoveByID(dao.BsonId(param.ID)); err != nil {
		return err
	}
	return nil
}

func (*_RuleManager) Get(param models.RuleQuery) (interface{}, error) {
	var (
		result models.RuleQueryResult
		resp   []*models.RuleInsert
		query  = bson.M{}
	)
	if param.Type != 0 {
		query["type"] = param.Type
	}
	if param.Status != "" {
		query["status"] = utils.StrToBool(param.Status)
	}
	if param.Search != "" {
		query["$or"] = []bson.M{
			{"name": bson.M{"$regex": param.Search, "$options": "$i"}},
			{"target_host": bson.M{"$regex": param.Search, "$options": "$i"}},
			{"target_port": bson.M{"$regex": param.Search, "$options": "$i"}},
		}
	}
	if err := repo.SelectWithPage(query, param.Page, param.Size, &resp, "-updated_time"); err != nil {
		return nil, err
	}
	result.Size = param.Size
	result.Page = param.Page
	result.Total = repo.Count(query)
	result.Items = models.RuleQueryResultFunc(resp)
	result.Pages = int(math.Ceil(float64(result.Total) / float64(param.Size)))
	return result, nil
}

func (*_RuleManager) Put(param models.QueryID, body models.Rule) error {
	var (
		r models.RuleInsert
	)
	if !bson.IsObjectIdHex(param.ID) {
		return fmt.Errorf("invalid ObjectIdHex")
	}
	if err := repo.SelectById(dao.BsonId(param.ID), &r); err != nil {
		return err
	}
	r.Rule = body
	r.UpdatedTime = time.Now().Unix()
	if err := repo.UpdateById(dao.BsonId(param.ID), &r); err != nil {
		return err
	}
	return nil
}

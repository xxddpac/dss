package management

import (
	"dss/core/dao"
	"dss/core/global"
	"dss/core/models"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"math"
	"time"
)

var (
	ScanItemManager *_ScanItemManager
)

type _ScanItemManager struct {
}

func (*_ScanItemManager) Get(param models.ScanItemQuery) (interface{}, error) {
	var (
		result models.ScanItemQueryResult
		resp   []*models.ScanItemInsert
		query  = bson.M{}
	)
	if param.Level != 0 {
		query["level"] = param.Level
	}
	if param.Search != "" {
		query["$or"] = []bson.M{
			{"name": bson.M{"$regex": param.Search, "$options": "$i"}},
			{"desc": bson.M{"$regex": param.Search, "$options": "$i"}},
		}
	}
	if err := dao.Repo(global.ScanItem).SelectWithPage(query, param.Page, param.Size, &resp, "-updated_time"); err != nil {
		return nil, err
	}
	result.Size = param.Size
	result.Page = param.Page
	result.Total = dao.Repo(global.ScanItem).Count(query)
	result.Items = models.ScanItemQueryResultFunc(resp)
	result.Pages = int(math.Ceil(float64(result.Total) / float64(param.Size)))
	return result, nil
}

func (*_ScanItemManager) Post(param models.ScanItem) error {
	if err := dao.Repo(global.ScanItem).Insert(models.ScanItemInsertFunc(param)); err != nil {
		return err
	}
	return nil
}

func (*_ScanItemManager) Put(param models.ScanItem, query models.QueryID) error {
	var (
		s models.ScanItemInsert
	)
	if !bson.IsObjectIdHex(query.ID) {
		return fmt.Errorf("invalid ObjectIdHex")
	}
	if err := dao.Repo(global.ScanItem).SelectById(dao.BsonId(query.ID), &s); err != nil {
		return err
	}
	s.ScanItem = param
	s.UpdatedTime = time.Now().Unix()
	if err := dao.Repo(global.ScanItem).UpdateById(dao.BsonId(query.ID), &s); err != nil {
		return err
	}
	return nil
}

func (*_ScanItemManager) Delete(query models.QueryID) error {
	if !bson.IsObjectIdHex(query.ID) {
		return fmt.Errorf("invalid ObjectIdHex")
	}
	if err := dao.Repo(global.ScanItem).RemoveByID(dao.BsonId(query.ID)); err != nil {
		return err
	}
	return nil
}

func (*_ScanItemManager) Query(query models.QueryID) (interface{}, error) {
	var (
		resp []*models.ScanItemInsert
	)
	if !bson.IsObjectIdHex(query.ID) {
		return nil, fmt.Errorf("invalid ObjectIdHex")
	}
	if err := dao.Repo(global.ScanItem).Select(bson.M{"_id": dao.BsonId(query.ID)}, &resp); err != nil {
		return nil, err
	}
	result := models.ScanItemQueryResultFunc(resp)
	if len(result) == 0 {
		return nil, fmt.Errorf("query id %v err", query.ID)
	}
	return result[0], nil
}

func (*_ScanItemManager) Enum() interface{} {
	var (
		res = make(map[string][]map[string]interface{})
	)
	res["level"] = append(res["level"],
		map[string]interface{}{"key": global.Critical.String(), "value": global.Critical},
		map[string]interface{}{"key": global.High.String(), "value": global.High},
		map[string]interface{}{"key": global.Middle.String(), "value": global.Middle},
		map[string]interface{}{"key": global.Low.String(), "value": global.Low},
	)
	return res
}

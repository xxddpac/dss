package management

import (
	"dss/common/log"
	"dss/common/utils"
	"dss/core/dao"
	"dss/core/models"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"math"
	"time"
)

var (
	PortManager *_PortManager
)

type _PortManager struct {
}

func (*_PortManager) Get(param models.ScanQuery) (interface{}, error) {
	var (
		result models.ScanQueryResult
		resp   []*models.ScanInsert
		query  = bson.M{}
		repo   = dao.Repository{Collection: "port_scan"}
	)
	if param.Location != "" {
		query["location"] = param.Location
	}
	if param.Date != "" {
		query["done_time"] = param.Date
	} else {
		query["done_time"] = time.Now().Format(utils.TimeLayout)
	}
	if param.Search != "" {
		query["$or"] = []bson.M{
			{"host": bson.M{"$regex": param.Search, "$options": "$i"}},
			{"port": bson.M{"$regex": param.Search, "$options": "$i"}},
		}
	}
	if err := repo.SelectWithPage(query, param.Page, param.Size, &resp, "-updated_time"); err != nil {
		return nil, err
	}
	if param.Date == "" {
		if len(resp) == 0 {
			query["done_time"] = time.Now().AddDate(0, 0, -1).Format(utils.TimeLayout)
			if err := repo.SelectWithPage(query, param.Page, param.Size, &resp, "-updated_time"); err != nil {
				return nil, err
			}
		}
	}
	result.Size = param.Size
	result.Page = param.Page
	result.Total = repo.Count(query)
	result.Items = models.ScanQueryResultFunc(resp)
	result.Pages = int(math.Ceil(float64(result.Total) / float64(param.Size)))
	return result, nil
}

func (*_PortManager) FieldGroupBy(field string) (pipeline []bson.M) {
	group := bson.M{"$group": bson.M{"_id": fmt.Sprintf("$%s", field), "count": bson.M{"$sum": 1}}}
	orderBy := bson.M{"$sort": bson.M{"count": 1}}
	pipeline = []bson.M{group, orderBy}
	return
}

func (*_PortManager) Location() (interface{}, error) {
	var (
		err            error
		field          = "location"
		resp, pipeline []bson.M
		result         []string
		repo           = dao.Repository{Collection: "port_scan"}
	)
	pipeline = PortManager.FieldGroupBy(field)
	if err = repo.Aggregate(pipeline, &resp); err != nil {
		return nil, err
	}
	for _, item := range resp {
		if _, ok := item["_id"].(string); !ok {
			continue
		}
		result = append(result, item["_id"].(string))
	}
	return result, nil
}

func (*_PortManager) Clear() {
	var (
		err            error
		field          = "done_time"
		result         []string
		resp, pipeline []bson.M
		repo           = dao.Repository{Collection: "port_scan"}
	)
	pipeline = PortManager.FieldGroupBy(field)
	if err = repo.Aggregate(pipeline, &resp); err != nil {
		log.Errorf("group by field err:%v", err)
		return
	}
	for _, item := range resp {
		if _, ok := item["_id"].(string); !ok {
			continue
		}
		result = append(result, item["_id"].(string))
	}
	if len(result) > 7 {
		result = result[:len(result)-7]
		for _, item := range result {
			if err = repo.RemoveAll(bson.M{field: item}); err != nil {
				log.Errorf("remove field err:%v", err)
			}
		}
	}
}

package management

import (
	"github.com/globalsign/mgo/bson"
	"goportscan/common/utils"
	"goportscan/core/dao"
	"goportscan/core/models"
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

package management

import (
	"github.com/globalsign/mgo/bson"
	"goportscan/core/dao"
	"goportscan/core/models"
	"math"
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
	if param.Search != "" {
		query["$or"] = []bson.M{
			{"host": bson.M{"$regex": param.Search, "$options": "$i"}},
			{"port": bson.M{"$regex": param.Search, "$options": "$i"}},
		}
	}
	if err := repo.SelectWithPage(query, param.Page, param.Size, &resp, "-updated_time"); err != nil {
		return nil, err
	}
	result.Size = param.Size
	result.Page = param.Page
	result.Total = repo.Count(query)
	result.Items = models.ScanQueryResultFunc(resp)
	result.Pages = int(math.Ceil(float64(result.Total) / float64(param.Size)))
	return result, nil
}

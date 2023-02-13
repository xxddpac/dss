package management

import (
	"dss/core/dao"
	"dss/core/global"
	"dss/core/models"
	"github.com/globalsign/mgo/bson"
	"math"
)

var (
	GrpcManager *_GrpcManager
)

type _GrpcManager struct {
}

func (*_GrpcManager) Get(param models.ClientQuery) (interface{}, error) {
	var (
		result models.ClientQueryResult
		resp   []*models.ClientInsert
		query  = bson.M{}
	)
	if param.Search != "" {
		query["$or"] = []bson.M{
			{"name": bson.M{"$regex": param.Search, "$options": "$i"}},
			{"host": bson.M{"$regex": param.Search, "$options": "$i"}},
			{"platform": bson.M{"$regex": param.Search, "$options": "$i"}},
			{"platform_version": bson.M{"$regex": param.Search, "$options": "$i"}},
		}
	}
	if err := dao.Repo(global.GrpcClient).SelectWithPage(query, param.Page, param.Size, &resp, "-updated_time"); err != nil {
		return nil, err
	}
	result.Size = param.Size
	result.Page = param.Page
	result.Total = dao.Repo(global.GrpcClient).Count(query)
	result.Items = models.ClientQueryResultFunc(resp)
	result.Pages = int(math.Ceil(float64(result.Total) / float64(param.Size)))
	return result, nil
}

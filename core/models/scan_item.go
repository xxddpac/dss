package models

import (
	"dss/core/global"
	"github.com/globalsign/mgo/bson"
	"time"
)

type ScanItem struct {
	Name  string       `json:"name" bson:"name" binding:"required"`
	Desc  string       `json:"desc" bson:"desc" binding:"required"`
	Level global.Level `json:"level" bson:"level" binding:"required"`
}

type ScanItemInsert struct {
	ScanItem `bson:",inline"`
	BasePo   `bson:",inline"`
}

func ScanItemInsertFunc(s ScanItem) *ScanItemInsert {
	return &ScanItemInsert{
		ScanItem: s,
		BasePo: BasePo{
			Id:          bson.NewObjectId(),
			CreatedTime: time.Now().Unix(),
			UpdatedTime: time.Now().Unix(),
		},
	}
}

type ScanItemQuery struct {
	QueryPage
	Search string       `form:"search"`
	Level  global.Level `form:"level"`
}

func ScanItemQueryFunc() *ScanItemQuery {
	return &ScanItemQuery{
		QueryPage: QueryPage{Page: 1, Size: 10},
	}
}

type ScanItemQueryResult struct {
	QueryResult
	Items []ScanItemQueryResultDto `json:"items"`
}

type ScanItemQueryResultDto struct {
	ScanItemQueryDto
	LevelDesc string `json:"level_desc" bson:"level_desc"`
}

type ScanItemQueryDto struct {
	BaseDto
	ScanItem
}

func ScanItemQueryResultFunc(s []*ScanItemInsert) []ScanItemQueryResultDto {
	var (
		resp   ScanItemQueryResultDto
		result []ScanItemQueryResultDto
	)
	for _, item := range s {
		resp.ScanItemQueryDto = *item.ToDto()
		resp.ScanItem = item.ScanItem
		resp.LevelDesc = item.Level.String()
		result = append(result, resp)
	}
	return result
}

func (s *ScanItemInsert) ToDto() *ScanItemQueryDto {
	dto := &ScanItemQueryDto{}
	dto.BaseDto = *(&s.BasePo).ToDto()
	return dto
}

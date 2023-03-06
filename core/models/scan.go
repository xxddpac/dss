package models

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type Scan struct {
	RuleId   string `json:"rule_id" bson:"rule_id"`
	TaskId   string `json:"task_id" bson:"task_id"`
	Host     string `json:"host" bson:"host"`
	Port     string `json:"port" bson:"port"`
	Location string `json:"location" bson:"location"`
	DoneTime string `json:"done_time" bson:"done_time"`
}

type ScanInsert struct {
	Scan   `bson:",inline"`
	BasePo `bson:",inline"`
}

func ScanInsertFunc(s Scan) *ScanInsert {
	return &ScanInsert{
		Scan: s,
		BasePo: BasePo{
			Id:          bson.NewObjectId(),
			CreatedTime: time.Now().Unix(),
			UpdatedTime: time.Now().Unix(),
		},
	}
}

type ScanQuery struct {
	QueryPage
	Date     string `form:"date"`
	Location string `form:"location"`
	Search   string `form:"search"`
}

func ScanQueryFunc() *ScanQuery {
	return &ScanQuery{
		QueryPage: QueryPage{Page: 1, Size: 10},
	}
}

type ScanQueryResult struct {
	QueryResult
	Items []ScanQueryResultDto `json:"items"`
}

type ScanQueryResultDto struct {
	ScanQueryDto
}

type ScanQueryDto struct {
	BaseDto
	Scan
}

func ScanQueryResultFunc(s []*ScanInsert) []ScanQueryResultDto {
	var (
		resp   ScanQueryResultDto
		result []ScanQueryResultDto
	)
	for _, item := range s {
		resp.ScanQueryDto = *item.ToDto()
		resp.Scan = item.Scan
		result = append(result, resp)
	}
	return result
}

func (s *ScanInsert) ToDto() *ScanQueryDto {
	dto := &ScanQueryDto{}
	dto.BaseDto = *(&s.BasePo).ToDto()
	return dto
}

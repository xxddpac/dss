package models

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type Scan struct {
	Host string `json:"host" bson:"host"`
	Port string `json:"port" bson:"port"`
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
	Search string `form:"search"`
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

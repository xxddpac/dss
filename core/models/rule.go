package models

import (
	"dss/core/global"
	"github.com/globalsign/mgo/bson"
	"time"
)

type Rule struct {
	Name       string          `json:"name" bson:"name" binding:"required"`
	Status     bool            `json:"status" bson:"status" binding:"required"`
	TargetHost string          `json:"target_host" bson:"target_host" binding:"required"`
	TargetPort string          `json:"target_port" bson:"target_port" binding:"required"`
	ScanItemId []string        `json:"scan_item_id" bson:"scan_item_id" binding:"required"`
	Type       global.RuleType `json:"type" bson:"type" binding:"required"`
	Location   string          `json:"location" bson:"location"`
}

type RuleInsert struct {
	Rule   `bson:",inline"`
	BasePo `bson:",inline"`
}

func RuleInsertFunc(r Rule) *RuleInsert {
	return &RuleInsert{
		Rule: r,
		BasePo: BasePo{
			Id:          bson.NewObjectId(),
			CreatedTime: time.Now().Unix(),
			UpdatedTime: time.Now().Unix(),
		},
	}
}

type RuleQuery struct {
	QueryPage
	Status string          `form:"status"`
	Type   global.RuleType `form:"type"`
	Search string          `form:"search"`
}

func RuleQueryFunc() *RuleQuery {
	return &RuleQuery{
		QueryPage: QueryPage{Page: 1, Size: 10},
	}
}

type RuleQueryResult struct {
	QueryResult
	Items []RuleQueryResultDto `json:"items"`
}

type RuleQueryResultDto struct {
	RuleQueryDto
	TypeDesc string `json:"type_desc"`
}

type RuleQueryDto struct {
	BaseDto
	Rule
}

func RuleQueryResultFunc(r []*RuleInsert) []RuleQueryResultDto {
	var (
		resp   RuleQueryResultDto
		result []RuleQueryResultDto
	)
	for _, item := range r {
		resp.RuleQueryDto = *item.ToDto()
		resp.Rule = item.Rule
		resp.TypeDesc = item.Type.String()
		result = append(result, resp)
	}
	return result
}

func (r *RuleInsert) ToDto() *RuleQueryDto {
	dto := &RuleQueryDto{}
	dto.BaseDto = *(&r.BasePo).ToDto()
	return dto
}

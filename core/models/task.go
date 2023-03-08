package models

import (
	"dss/core/global"
	"github.com/globalsign/mgo/bson"
	"time"
)

type Task struct {
	RuleId       string             `json:"rule_id" bson:"rule_id"`
	Name         string             `json:"name" bson:"name"`
	Status       global.TaskStatus  `json:"status" bson:"status"`
	Count        int                `json:"count" bson:"count"`
	Progress     string             `json:"progress" bson:"progress"`
	ExecutedTime string             `json:"executed_time" bson:"executed_time"`
	RunType      global.TaskRunType `json:"run_type" bson:"run_type"`
}

type TaskInsert struct {
	Task   `bson:",inline"`
	BasePo `bson:",inline"`
}

func TaskInsertFunc(t Task) *TaskInsert {
	return &TaskInsert{
		Task: t,
		BasePo: BasePo{
			Id:          bson.NewObjectId(),
			CreatedTime: time.Now().Unix(),
			UpdatedTime: time.Now().Unix(),
		},
	}
}

type TaskQuery struct {
	QueryPage
	Status global.TaskStatus `form:"status"`
	Search string            `form:"search"`
}

func TaskQueryFunc() *TaskQuery {
	return &TaskQuery{
		QueryPage: QueryPage{Page: 1, Size: 10},
	}
}

type TaskQueryResult struct {
	QueryResult
	Items []TaskQueryResultDto `json:"items"`
}

type TaskQueryResultDto struct {
	TaskQueryDto
	StatusDesc string `json:"status_desc"`
}

type TaskQueryDto struct {
	BaseDto
	Task
}

func TaskQueryResultFunc(t []*TaskInsert) []TaskQueryResultDto {
	var (
		resp   TaskQueryResultDto
		result []TaskQueryResultDto
	)
	for _, item := range t {
		resp.TaskQueryDto = *item.ToDto()
		resp.Task = item.Task
		resp.StatusDesc = item.Status.String()
		result = append(result, resp)
	}
	return result
}

func (t *TaskInsert) ToDto() *TaskQueryDto {
	dto := &TaskQueryDto{}
	dto.BaseDto = *(&t.BasePo).ToDto()
	return dto
}

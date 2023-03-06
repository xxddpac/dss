package models

import (
	"dss/core/global"
	"github.com/globalsign/mgo/bson"
	"time"
)

type Task struct {
	RuleId       string            `json:"rule_id" bson:"rule_id"`
	Name         string            `json:"name" bson:"name"`
	Status       global.TaskStatus `json:"status" bson:"status"`
	Count        int               `json:"count" bson:"count"`
	Progress     string            `json:"progress" bson:"progress"`
	ExecutedTime string            `json:"executed_time" bson:"executed_time"`
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

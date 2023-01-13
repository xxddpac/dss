package models

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type Rule struct {
	Name       string `json:"name" bson:"name" binding:"required"`
	Status     bool   `json:"status" bson:"status" binding:"required"`
	TargetHost string `json:"target_host" bson:"target_host" binding:"required"`
	TargetPort string `json:"target_port" bson:"target_port" binding:"required"`
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

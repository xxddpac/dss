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

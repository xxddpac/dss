package models

import (
	"github.com/globalsign/mgo/bson"
)

type BasePo struct {
	Id          bson.ObjectId `bson:"_id"`
	CreatedTime int64         `bson:"created_time,omitempty"`
	UpdatedTime int64         `bson:"updated_time,omitempty"`
}

type QueryID struct {
	ID string `form:"id" binding:"required"`
}

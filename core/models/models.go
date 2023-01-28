package models

import (
	"dss/common/utils"
	"github.com/globalsign/mgo/bson"
)

type BasePo struct {
	Id          bson.ObjectId `bson:"_id"`
	CreatedTime int64         `bson:"created_time,omitempty"`
	UpdatedTime int64         `bson:"updated_time,omitempty"`
}

type BaseDto struct {
	Id          string `json:"id"`
	CreatedTime string `json:"created_time,omitempty"`
	UpdatedTime string `json:"updated_time,omitempty"`
}

func (b *BasePo) ToDto() *BaseDto {
	dto := &BaseDto{}
	dto.Id = b.Id.Hex()
	dto.CreatedTime = utils.UnixToString(b.CreatedTime)
	dto.UpdatedTime = utils.UnixToString(b.UpdatedTime)

	return dto
}

type QueryID struct {
	ID string `form:"id" binding:"required"`
}

type QueryPage struct {
	Page int `form:"page" binding:"gte=1"`
	Size int `form:"size" binding:"gte=1"`
}

type QueryResult struct {
	Page  int `json:"page"`
	Size  int `json:"size"`
	Total int `json:"total"`
	Pages int `json:"pages"`
}

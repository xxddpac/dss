package models

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type Client struct {
	Host            string `json:"host" bson:"host"`
	Name            string `json:"name" bson:"name"`
	Platform        string `json:"platform" bson:"platform"`
	PlatformVersion string `json:"platform_version" bson:"platform_version"`
}

type ClientInsert struct {
	Client `bson:",inline"`
	BasePo `bson:",inline"`
}

func ClientInsertFunc(c Client) *ClientInsert {
	return &ClientInsert{
		Client: c,
		BasePo: BasePo{
			Id:          bson.NewObjectId(),
			CreatedTime: time.Now().Unix(),
			UpdatedTime: time.Now().Unix(),
		},
	}
}

type ClientQuery struct {
	QueryPage
	Search string `form:"search"`
}

func ClientQueryFunc() *ClientQuery {
	return &ClientQuery{
		QueryPage: QueryPage{Page: 1, Size: 10},
	}
}

type ClientQueryResult struct {
	QueryResult
	Items []ClientQueryResultDto `json:"items"`
}

type ClientQueryResultDto struct {
	ClientQueryDto
	IsOnline bool `json:"is_online"`
}

type ClientQueryDto struct {
	BaseDto
	Client
}

func ClientQueryResultFunc(c []*ClientInsert) []ClientQueryResultDto {
	var (
		resp   ClientQueryResultDto
		result []ClientQueryResultDto
	)
	for _, item := range c {
		resp.ClientQueryDto = *item.ToDto()
		resp.Client = item.Client
		result = append(result, resp)
	}
	return result
}

func (c *ClientInsert) ToDto() *ClientQueryDto {
	dto := &ClientQueryDto{}
	dto.BaseDto = *(&c.BasePo).ToDto()
	return dto
}

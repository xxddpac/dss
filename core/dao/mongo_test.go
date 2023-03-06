package dao

import (
	"dss/common/mongo"
	"dss/core/global"
	"dss/core/models"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"testing"
)

var (
	fakeMongoConfig = &mongo.Config{
		Host:     "10.101.191.106:27017",
		Database: "SecurityManagement",
		Auth: &mongo.AuthConfig{
			User:     "admin",
			Passwd:   "123456",
			Database: "admin",
		},
	}
)

func TestRepository_BulkWrite(t *testing.T) {
	if err := mongo.Init(fakeMongoConfig); err != nil {
		t.Fatal(err)
	}
	var resp []interface{}
	resp = append(resp,
		models.ScanInsertFunc(models.Scan{
			Host: "1.1.1.1",
			Port: "22",
		}),
		models.ScanInsertFunc(models.Scan{
			Host: "1.1.1.2",
			Port: "23",
		}))
	if err := Repo(global.Scan).BulkWrite(resp); err != nil {
		t.Fatal(err)
	}
}

func TestRepository_Insert(t *testing.T) {
	if err := mongo.Init(fakeMongoConfig); err != nil {
		t.Fatal(err)
	}
	if err := Repo(global.ScanRule).Insert(models.RuleInsertFunc(models.Rule{
		Name:       "test",
		Status:     true,
		TargetHost: "1.1.1.0/24",
		TargetPort: "1-65535",
	})); err != nil {
		t.Fatal(err)
	}
}

func TestRepository_Select(t *testing.T) {
	if err := mongo.Init(fakeMongoConfig); err != nil {
		t.Fatal(err)
	}
	var (
		ruleSlice []models.Rule
	)
	if err := Repo(global.ScanRule).SelectAll(&ruleSlice); err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(ruleSlice))
	for _, item := range ruleSlice {
		fmt.Println(item)
	}
}

func TestRepository_Aggregate(t *testing.T) {
	var (
		resp []bson.M
	)
	if err := mongo.Init(fakeMongoConfig); err != nil {
		t.Fatal(err)
	}
	group := bson.M{"$group": bson.M{"_id": "$location", "count": bson.M{"$sum": 1}}}
	orderBy := bson.M{"$sort": bson.M{"count": 1}}
	pipeline := []bson.M{group, orderBy}
	if err := Repo(global.Scan).Aggregate(pipeline, &resp); err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp)
}

package dao

import (
	"goportscan/common/mongo"
	"goportscan/core/models"
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

func TestMongo(t *testing.T) {
	if err := mongo.Init(fakeMongoConfig); err != nil {
		t.Fatal(err)
	}
	var resp []interface{}
	repo := &Repository{"port_scan"}
	resp = append(resp,
		models.ScanInsertFunc(models.Scan{
			Host: "1.1.1.1",
			Port: "22",
		}),
		models.ScanInsertFunc(models.Scan{
			Host: "1.1.1.2",
			Port: "23",
		}))
	if err := repo.BulkWrite(resp); err != nil {
		t.Fatal(err)
	}
}

package dao

import (
	"goportscan/common/mongo"
)

type Repository struct {
	Collection string
}

func (r *Repository) BulkWrite(docs []interface{}) error {
	client := mongo.GetConn(r.Collection)
	defer client.Close()
	col := client.Collection().Bulk()
	for _, doc := range docs {
		col.Insert(doc)
	}
	_, err := col.Run()
	return err
}

func (r *Repository) Insert(doc interface{}) error {
	client := mongo.GetConn(r.Collection)
	defer client.Close()
	return client.Collection().Insert(doc)
}

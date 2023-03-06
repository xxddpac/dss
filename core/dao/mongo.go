package dao

import (
	"dss/common/mongo"
	"github.com/globalsign/mgo/bson"
	"time"
)

type Repository struct {
	Collection string
}

func BsonId(id string) bson.ObjectId {
	return bson.ObjectIdHex(id)
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

func (r *Repository) SelectAll(result interface{}, fields ...string) error {
	client := mongo.GetConn(r.Collection)
	defer client.Close()
	if err := client.Collection().Find(nil).Sort(fields...).All(result); nil != err {
		return err
	}
	return nil
}

func (r *Repository) RemoveByID(id interface{}) error {
	client := mongo.GetConn(r.Collection)
	defer client.Close()
	return client.Collection().RemoveId(id)
}

func (r *Repository) Count(query bson.M) int {
	client := mongo.GetConn(r.Collection)
	defer client.Close()
	num, _ := client.Collection().Find(query).Count()
	return num
}

func (r *Repository) SelectWithPage(query bson.M, page, size int, result interface{}, fields ...string) error {
	client := mongo.GetConn(r.Collection)
	defer client.Close()
	limit := size
	skip := (page - 1) * size
	if err := client.Collection().Find(query).Sort(fields...).Skip(skip).Limit(limit).All(result); nil != err {
		return err
	}
	return nil
}

func (r *Repository) SelectById(id interface{}, result interface{}) error {
	client := mongo.GetConn(r.Collection)
	defer client.Close()
	if err := client.Collection().FindId(id).One(result); nil != err {
		return err
	}
	return nil
}

func (r *Repository) UpdateById(id interface{}, update interface{}) error {
	client := mongo.GetConn(r.Collection)
	defer client.Close()
	return client.Collection().UpdateId(id, update)
}

func (r *Repository) Aggregate(query []bson.M, resp interface{}) error {
	client := mongo.GetConn(r.Collection)
	defer client.Close()
	if err := client.Collection().Pipe(query).All(resp); err != nil {
		return err
	}
	return nil
}

func (r *Repository) RemoveAll(query bson.M) error {
	client := mongo.GetConn(r.Collection)
	defer client.Close()
	_, err := client.Collection().RemoveAll(query)
	return err
}

func (r *Repository) Select(query bson.M, result interface{}) error {
	client := mongo.GetConn(r.Collection)
	defer client.Close()
	if err := client.Collection().Find(query).All(result); err != nil {
		return err
	}
	return nil
}

func (r *Repository) SetField(id interface{}, update bson.M) error {
	update["updated_time"] = time.Now().Unix()
	client := mongo.GetConn(r.Collection)
	defer client.Close()
	return client.Collection().Update(bson.M{"_id": id}, bson.M{"$set": update})
}

package achieve

import (
	"go-vrs/consume"
	"go-vrs/lib/mong"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoAch struct {
	Mgo *mong.Mgo
}

type ResultAch struct {
	res *mongo.InsertManyResult
}

type ResultOneAch struct {
	res *mongo.InsertOneResult
}

func (m *MongoAch) Send(d consume.Data) (consume.Result, error) {
	data := d.([]interface{})
	if len(data) == 0 {
		return nil, nil
	}
	res, err := m.Mgo.InsertMany(data)
	return ResultAch{res}, err
}
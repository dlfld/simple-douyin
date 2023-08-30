package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	_ "go.mongodb.org/mongo-driver/mongo/readpref"
	"sync"
)

var mdb *mongo.Client
var once sync.Once
var err error

func NewMongodbConn() (*mongo.Client, error) {
	once.Do(func() {
		mdb, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://42.192.46.30:27017"))
	})
	if err != nil {
		return nil, err
	}
	return mdb, nil
}

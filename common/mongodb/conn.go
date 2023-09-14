package mongodb

import (
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	_ "go.mongodb.org/mongo-driver/mongo/readpref"
)

var mdb *mongo.Client
var once sync.Once
var err error

func NewMongodbConn() (*mongo.Client, error) {
	once.Do(func() {
		// Set client options
		credential := options.Credential{
			Username: "admin",
			Password: "abc123456",
		}
		clientOptions := options.Client().ApplyURI("mongodb://mongodb:27017").SetAuth(credential)
		mdb, err = mongo.Connect(context.Background(), clientOptions)
	})
	if err != nil {
		return nil, err
	}
	return mdb, nil
}

package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"testing"
)

func TestNewMongodbConn(t *testing.T) {
	mdb, _ := NewMongodbConn()
	if err := mdb.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
}

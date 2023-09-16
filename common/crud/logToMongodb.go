package crud

import (
	"context"
	"fmt"

	"github.com/douyin/common/kafkaLog"
	"github.com/douyin/common/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

var logType = [7]string{"Trace", "Debug", "Info", "Notice", "Warn", "Error", "Fatal"}

func SetLog(serviceName string, log *kafkaLog.LogRecord) {
	mdb, _ := mongodb.NewMongodbConn()
	collection := mdb.Database("log").Collection(serviceName)

	mlog := bson.D{{"Type", logType[log.Type]}, {"Msg", log.Value}, {"Time", log.Time}}
	result, err := collection.InsertOne(context.TODO(), mlog)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.InsertedID)
}

// GetLog get all log of the service
func GetLog(serviceName string) {
	mdb, _ := mongodb.NewMongodbConn()
	collection := mdb.Database("log").Collection(serviceName)

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	fmt.Println("displaying all results in a collection")
	for _, result := range results {
		fmt.Println(result)
	}
}

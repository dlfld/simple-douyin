package productor

import (
	"fmt"
	"github.com/douyin/common/conf"
	"testing"
	"time"
)

func TestWriteLogToKafka(t *testing.T) {
	n := 0
	for ; n < 1000; n++ {
		fmt.Println("======")
		collector, err := NewLogCollector(conf.MessageService.Name)
		if err != nil {
			panic(err)
		}

		collector.Info(fmt.Sprintf("%d: 这是一条测试消息", time.Now().Unix()))
		time.Sleep(1 * time.Second)
	}

}

//docker run -d --network=kafka_network --name zookeeper -p 2181:2181 -e ZOOKEEPER_CLIENT_PORT=2181 wurstmeister/zookeeper
//
//docker run -d --network=kafka_network --name kafka -p 9092:9092 -e KAFKA_ADVERTISED_HOST_NAME=kafka -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://124.223.117.87:9092 -e KAFKA_LISTENERS=PLAINTEXT://0.0.0.0:9092 -e KAFKA_ZOOKEEPER_CONNECT=zookeeper:2181 wurstmeister/kafka

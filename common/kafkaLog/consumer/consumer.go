package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/douyin/common/kafkaLog"
	"github.com/segmentio/kafka-go"
	"time"
)

var (
	kr      *kafka.Reader
	groupId = "loggerHandler1"
)

func init() {
	kr = kafka.NewReader(
		kafka.ReaderConfig{
			Brokers:        []string{fmt.Sprintf(kafkaLog.KafkaAddr)},
			GroupID:        groupId,
			Topic:          kafkaLog.Topic,
			CommitInterval: time.Second * 2,
			StartOffset:    kafka.FirstOffset,
			MinBytes:       1024,
			MaxBytes:       1024 * 1024 * 1024,
		},
	)
}

// ReadLogFromKafka 从kafka中读取日志信息
func ReadLogFromKafka() (serviceName string, log *kafkaLog.LogRecord, err error) {
	msg, err := kr.ReadMessage(context.Background())
	if err != nil {
		kafkaLog.KafkaLogger.Errorf("failed read kafka, err=%s", err.Error())
		return "", nil, err
	}
	log = new(kafkaLog.LogRecord)
	if err = json.Unmarshal(msg.Value, log); err != nil {
		panic(err)
	}
	return string(msg.Key), log, nil
}

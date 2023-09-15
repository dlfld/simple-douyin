package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/douyin/common/conf"
	"github.com/douyin/common/kafkaLog"
	"github.com/segmentio/kafka-go"
	"time"
)

var (
	logKr         *kafka.Reader
	messageKr     *kafka.Reader
	relationKr    *kafka.Reader
	videoKr       *kafka.Reader
	interactionKr *kafka.Reader
	userKr        *kafka.Reader
)

func init() {
	logKr = newKafkaReader(kafkaLog.Topic, kafkaLog.Topic+"-group")
	messageKr = newKafkaReader(conf.MessageService.Name, conf.MessageService.Name+"-group")
	relationKr = newKafkaReader(conf.RelationService.Name, conf.RelationService.Name+"-group")
	videoKr = newKafkaReader(conf.VideoService.Name, conf.VideoService.Name+"-group")
	interactionKr = newKafkaReader(conf.InteractionService.Name, conf.InteractionService.Name+"-group")
	userKr = newKafkaReader(conf.UserService.Name, conf.UserService.Name+"-group")
}

func newKafkaReader(topic, groupId string) *kafka.Reader {
	return kafka.NewReader(
		kafka.ReaderConfig{
			Brokers:        []string{fmt.Sprintf(kafkaLog.KafkaAddr)},
			GroupID:        groupId,
			Topic:          topic,
			CommitInterval: time.Second * 2,
			StartOffset:    kafka.FirstOffset,
			MinBytes:       1024,
			MaxBytes:       1024 * 1024 * 1024,
		})
}

// PopLog 从kafka中读取日志信息
func PopLog() (serviceName string, log *kafkaLog.LogRecord, err error) {
	msg, err := logKr.ReadMessage(context.Background())
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

func PopMessage() (key string, values []byte, err error) {
	return pop(messageKr)
}
func PopUser() (key string, values []byte, err error) {
	return pop(userKr)
}
func PopRelation() (key string, values []byte, err error) {
	return pop(relationKr)
}
func PopVideo() (key string, values []byte, err error) {
	return pop(videoKr)
}
func PopInteraction() (key string, values []byte, err error) {
	return pop(interactionKr)
}
func pop(kr *kafka.Reader) (key string, values []byte, err error) {
	msg, err := kr.ReadMessage(context.Background())
	if err != nil {
		kafkaLog.KafkaLogger.Errorf("failed read kafka, err=%s", err.Error())
		return
	}
	return string(msg.Key), msg.Value, nil
}

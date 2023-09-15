package productor

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/bytedance/sonic"
	"github.com/douyin/common/conf"
	"github.com/douyin/common/kafkaLog"
	"github.com/segmentio/kafka-go"
)

var (
	logKw         *kafka.Writer
	userKw        *kafka.Writer
	messageKw     *kafka.Writer
	relationKw    *kafka.Writer
	videoKw       *kafka.Writer
	interactionKw *kafka.Writer

	servicesMap map[string]struct{}
)

const (
	retryTime       = 2
	retryWaiterTime = 2
)

func init() {
	// 初始化kafka
	logKw = newKafkaWriter(kafkaLog.Topic)
	userKw = newKafkaWriter(conf.UserService.Name)
	messageKw = newKafkaWriter(conf.MessageService.Name)
	relationKw = newKafkaWriter(conf.RelationService.Name)
	videoKw = newKafkaWriter(conf.VideoService.Name)
	interactionKw = newKafkaWriter(conf.InteractionService.Name)
	// 记录所有服务名
	servicesMap = make(map[string]struct{})
	for _, serviceName := range conf.GetAllServiceName() {
		servicesMap[serviceName] = struct{}{}
	}
}
func newKafkaWriter(topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:                   kafka.TCP(kafkaLog.KafkaAddr),
		Balancer:               &kafka.Hash{},
		Topic:                  topic,
		AllowAutoTopicCreation: true,
		RequiredAcks:           kafka.RequireNone,
	}
}

type LogCollector struct {
	ServiceName string
}

func (l *LogCollector) Debug(logValue string) {
	pushLogToKafka(l.ServiceName, logger.LevelDebug, logValue)
}
func (l *LogCollector) Info(logValue string) {
	pushLogToKafka(l.ServiceName, logger.LevelInfo, logValue)
}
func (l *LogCollector) Waring(logValue string) {
	pushLogToKafka(l.ServiceName, logger.LevelWarn, logValue)
}
func (l *LogCollector) Error(logValue string) {
	pushLogToKafka(l.ServiceName, logger.LevelError, logValue)
}

// NewLogCollector new 日志收集器，通过收集器的方法可以异步的将日志写入kafka
// serviceName: 是当前微服务的名称，需要和conf保持一致
func NewLogCollector(serviceName string) (*LogCollector, error) {
	if _, ok := servicesMap[serviceName]; !ok {
		return nil, errors.New(fmt.Sprintf("服务名[%s]非法", serviceName))
	}
	return &LogCollector{
		ServiceName: serviceName,
	}, nil
}

// pushLogToKafka 向kafka写入消息，错误信息用key-value对表示
// key：服务名称 （和配置文件一致，服务名不一致将会拒绝写入）
// value: 服务日志信息
// 验证服务名合法后，就会开启协程去执行写入消息的操作
func pushLogToKafka(key string, level logger.Level, logValue string) {
	record, _ := sonic.Marshal(kafkaLog.LogRecord{
		Type:  level,
		Value: logValue,
	})
	go push(logKw, retryTime, key, record)
}

func PushMessageToKafka(key string, value []byte) {
	push(messageKw, retryTime, key, value)
}
func PushUserToKafka(key string, value []byte) {
	push(userKw, retryTime, key, value)
}
func PushRelationToKafka(key string, value []byte) {
	push(relationKw, retryTime, key, value)
}
func PushVideoToKafka(key string, value []byte) {
	push(videoKw, retryTime, key, value)
}
func PushInteractionToKafka(key string, value []byte) {
	push(interactionKw, retryTime, key, value)
}

func push(kw *kafka.Writer, reTime int, key string, value []byte) {
	err := kw.WriteMessages(
		context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: value,
		},
	)
	fmt.Printf("写入了一条消息service[%s], value[%s], top[%s]\n", key, value, logKw.Topic)
	// 写入失败，过一段时间后再重试
	if err != nil {
		// 没重试次数了，退出
		if reTime <= 0 {
			if err != nil {
				kafkaLog.KafkaLogger.Errorf("failed to push messages, key[%s], value[%s], err=%s:", key, value, err.Error())
				panic(err)
			}
			return
		}
		time.Sleep(time.Second * retryWaiterTime)
		push(kw, reTime-1, key, value)
	}
}

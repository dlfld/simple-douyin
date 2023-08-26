package productor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/douyin/common/conf"
	"github.com/douyin/common/kafkaLog"
	"github.com/segmentio/kafka-go"
	"time"
)

var (
	kw          *kafka.Writer
	servicesMap map[string]struct{}
)

const (
	retryTime       = 2
	retryWaiterTime = 2
)

func init() {
	// 初始化kafka
	kw = &kafka.Writer{
		Addr:                   kafka.TCP(kafkaLog.KafkaAddr),
		Balancer:               &kafka.Hash{},
		Topic:                  kafkaLog.Topic,
		AllowAutoTopicCreation: true,
		RequiredAcks:           kafka.RequireNone,
	}
	// 记录所有服务名
	servicesMap = make(map[string]struct{})
	for _, serviceName := range conf.GetAllServiceName() {
		servicesMap[serviceName] = struct{}{}
	}
}

type LogCollector struct {
	ServiceName string
}

func (l *LogCollector) Debug(logValue string) {
	writeLogToKafka(l.ServiceName, logger.LevelDebug, logValue)
}
func (l *LogCollector) Info(logValue string) {
	writeLogToKafka(l.ServiceName, logger.LevelInfo, logValue)
}
func (l *LogCollector) Waring(logValue string) {
	writeLogToKafka(l.ServiceName, logger.LevelWarn, logValue)
}
func (l *LogCollector) Error(logValue string) {
	writeLogToKafka(l.ServiceName, logger.LevelError, logValue)
}

func NewLogCollector(serviceName string) (*LogCollector, error) {
	if _, ok := servicesMap[serviceName]; !ok {
		return nil, errors.New(fmt.Sprintf("服务名[%s]非法", serviceName))
	}
	return &LogCollector{
		ServiceName: serviceName,
	}, nil
}

// writeLogToKafka 向kafka写入消息，错误信息用key-value对表示
// key：服务名称 （和配置文件一致，服务名不一致将会拒绝写入）
// value: 服务日志信息
// 验证服务名合法后，就会开启携程去执行写入消息的操作
func writeLogToKafka(key string, level logger.Level, logValue string) {
	record, _ := json.Marshal(kafkaLog.LogRecord{
		Type:  level,
		Value: logValue,
	})
	go write(retryTime, key, record)
}

func write(reTime int, key string, value []byte) {
	err := kw.WriteMessages(
		context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: value,
		},
	)
	fmt.Printf("写入了一条消息service[%s], value[%s], top[%s]\n", key, value, kw.Topic)
	// 写入失败，过一段时间后再重试
	if err != nil {
		// 没重试次数了，退出
		if reTime <= 0 {
			if err != nil {
				kafkaLog.KafkaLogger.Errorf("failed to write messages, key[%s], value[%s], err=%s:", key, value, err.Error())
				panic(err)
			}
			return
		}
		time.Sleep(time.Second * retryWaiterTime)
		write(reTime-1, key, value)
	}
}

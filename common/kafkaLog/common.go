package kafkaLog

import (
	"fmt"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/douyin/common/comLogger"
	"github.com/douyin/common/conf"
	"github.com/sirupsen/logrus"
)

var (
	KafkaLogger  *logrus.Logger
	KafkaLogFile = "./kafka.log"
	Topic        = "logger"
	KafkaAddr    = fmt.Sprintf("%s:%d", conf.Kafka.Addr, conf.Kafka.Port)
)

// LogRecord 一条日志记录
type LogRecord struct {
	Type  logger.Level // 日志类型
	Value string       // 日志信息
	Time  string       // 日志时间
}

func init() {
	// 初始化记录kafka错误信息的日志
	KafkaLogger = comLogger.NewLogger(KafkaLogFile)
}

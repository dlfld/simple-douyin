package comLogger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var logger *logrus.Logger

func NewLogger() *logrus.Logger {
	if logger == nil {
		logger = logrus.New()
		logger.SetReportCaller(true)
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,
		})
		file, err := os.OpenFile("./logfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			logger.SetOutput(file)
		} else {
			logger.Info("无法打开日志文件：", err)
		}
	}
	return logger
}

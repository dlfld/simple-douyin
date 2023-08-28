package main

import (
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/douyin/common/comLogger"
	"github.com/douyin/common/conf"
	"github.com/douyin/common/crud"
	"github.com/douyin/common/kafkaLog/consumer"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var loggerMap map[string]*logrus.Logger

func init() {
	loggerMap = make(map[string]*logrus.Logger, 0)
	for _, serviceName := range conf.GetAllServiceName() {
		loggerMap[serviceName] = comLogger.NewLogger(fmt.Sprintf("%s.log", serviceName))
	}
}

func handler() {
	fmt.Println("准备读取消息：")
	for {
		serviceName, log, err := consumer.ReadLogFromKafka()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		//set to mongodb
		crud.SetLog(serviceName, log)
		fmt.Println(serviceName, log)

		switch log.Type {
		case logger.LevelDebug:
			loggerMap[serviceName].Debug(log.Value)
		case logger.LevelInfo:
			loggerMap[serviceName].Info(log.Value)
		case logger.LevelWarn:
			loggerMap[serviceName].Warn(log.Value)
		case logger.LevelError:
			loggerMap[serviceName].Error(log.Value)
		}

		//fmt.Println("-------------------------------------")
		//fmt.Println("-------------------------------------")
		//fmt.Println("-------------------------------------")
		//crud.GetLog(serviceName)
	}
}

// 使用协程处理数据，当收到control c的退出信号或者程序报错时，等待5秒再退出，以让缓存中数据完全的写入到日志文件中
func main() {
	var sigs = make(chan os.Signal, 10)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		handler()
		sigs <- syscall.SIGTERM
	}()
	<-sigs
	fmt.Println("收到退出信号")
	time.Sleep(time.Second * 5)

}

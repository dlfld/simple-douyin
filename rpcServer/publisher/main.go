package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/douyin/common/kafkaLog/productor"
	videoconsumer "github.com/douyin/common/kafkaUpload/consumer"
)

var logger *productor.LogCollector

func init() {
	var err error
	if logger, err = productor.NewLogCollector("videoPublisher"); err != nil {
		panic(err)
	}

}
func main() {
	var sigs = make(chan os.Signal, 10)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for i := 0; i < 3; i++ {
			if err := videoconsumer.HandleVideoFromKafka(); err != nil {
				logger.Error(fmt.Sprintf("video pub failed, err=%s", err.Error()))
			}
			logger.Error(fmt.Sprintf("video pub restart %d", i))
		}
		sigs <- syscall.SIGTERM
	}()
	<-sigs
}

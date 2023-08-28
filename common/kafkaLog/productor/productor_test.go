package productor

import (
	"fmt"
	"github.com/douyin/common/conf"
	"testing"
	"time"
)

func TestWriteLogToKafka(t *testing.T) {
	for {
		collector, err := NewLogCollector(conf.MessageService.Name)
		if err != nil {
			panic(err)
		}
		collector.Error(fmt.Sprintf("%d: 这是一条测试消息", time.Now().Unix()))
	}

}

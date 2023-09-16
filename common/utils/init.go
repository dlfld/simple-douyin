package utils

import (
	"github.com/douyin/common/bloom"
	"github.com/douyin/common/conf"
	"github.com/douyin/common/kafkaLog/productor"
)

var LogCollector *productor.LogCollector
var bf *bloom.Filter

func init() {
	var err error
	//初始化日志收集器
	if LogCollector, err = productor.NewLogCollector(conf.VideoService.Name); err != nil {
		panic(err)
	}
	bf = bloom.NewBloom()
}

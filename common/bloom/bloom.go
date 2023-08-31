package bloom

import (
	"github.com/douyin/common/conf"
	"github.com/willf/bloom"
)

type Filter struct {
	filter *bloom.BloomFilter
}

func NewBloomFilter() *Filter {
	// 创建一个布隆过滤器，设置位数组大小为 1000000，使用 3 个哈希函数
	//filter := bloom.New(1000000, 3)
	filter := bloom.New(conf.BloomConfig.BloomBit, conf.BloomConfig.HashNum)
	return &Filter{
		filter: filter,
	}
}

const (
	videoPrefix   = "videoId_"
	userPrefix    = "userId_"
	commentPrefix = "commentId_"
)

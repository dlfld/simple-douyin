package bloom

import (
	"github.com/douyin/common/conf"
	"github.com/go-redis/redis/v8"
	"sync"
)

var once sync.Once

type Filter struct {
	filter *redis.Client
}

func NewBloom() *Filter {
	var rdb *redis.Client
	once.Do(func() {
		rdb = redis.NewClient(&redis.Options{Addr: conf.BloomConfig.Addr, Password: conf.BloomConfig.Password, DB: 0})
	})
	return &Filter{
		filter: rdb,
	}
}

func newBloomForTest() *Filter {
	var rdb *redis.Client
	once.Do(func() {
		rdb = redis.NewClient(&redis.Options{Addr: "localhost:6380", Password: "", DB: 0})
	})
	return &Filter{
		filter: rdb,
	}
}

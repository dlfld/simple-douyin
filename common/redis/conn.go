package redis

import (
	"context"
	"github.com/douyin/common/conf"
	"github.com/go-redis/redis/v8"
	"sync"
)

var rdb *redis.Client
var once sync.Once
var err error

func NewRedisConn() (*redis.Client, error) {
	once.Do(func() {
		rdb = redis.NewClient(&redis.Options{Addr: conf.Redis.Addr, Password: conf.Redis.Password, DB: 0})
		_, err = rdb.Ping(context.Background()).Result()
	})
	return rdb, err
}

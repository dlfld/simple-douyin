package dao

import (
	comrd "github.com/douyin/common/redis"

	grd "github.com/go-redis/redis/v8"
)

type redis struct {
	cli *grd.Client
}

func NewRedis() *redis {
	conn, _ := comrd.NewRedisConn()
	return &redis{cli: conn}
}

func (redis *redis) Test() {
	_ = redis.cli.String()
}

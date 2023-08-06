package main

import (
	"github.com/douyin/common/mysql"
	rdb "github.com/douyin/common/redis"
	"github.com/douyin/models"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var db *gorm.DB
var cache *redis.Client
var err error

func init() {
	db, err = mysql.NewMysqlConn()
	if err != nil {
		panic(ErrMysqlConn)
	}
	if err := db.AutoMigrate(&models.Message{}); err != nil {
		panic(err)
	}
	cache, err = rdb.NewRedisConn()
	if err != nil {
		panic(ErrCacheConn)
	}
}

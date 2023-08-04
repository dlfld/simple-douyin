package crud

import (
	"sync"

	"github.com/douyin/common/mysql"
	"github.com/douyin/common/oss"
	myredis "github.com/douyin/common/redis"
	"github.com/go-redis/redis/v8"

	"gorm.io/gorm"
)

var once sync.Once
var g *CachedCRUD

type CachedCRUD struct {
	redis *redis.Client
	mysql *gorm.DB
	oss   *oss.Service
}

func NewCachedCRUD() (*CachedCRUD, error) {
	var e error
	once.Do(
		func() {
			g = new(CachedCRUD)
			redis, e := myredis.NewRedisConn()
			if e != nil {
				return
			}
			g.redis = redis
			mysql, e := mysql.NewMysqlConn()
			if e != nil {
				return
			}
			g.mysql = mysql
			oss, e := oss.GetOssService()
			if e != nil {
				return
			}
			g.oss = oss
		},
	)

	return g, e
}

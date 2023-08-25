package crud

import (
	"sync"

	"github.com/douyin/common/mysql"
	"github.com/douyin/common/oss"
	myredis "github.com/douyin/common/redis"
	"github.com/go-redis/redis/v8"

	"gorm.io/gorm"
)

var crud *CachedCRUD
var once sync.Once

type CachedCRUD struct {
	redis *redis.Client
	mysql *gorm.DB
	oss   *oss.Service
}

func NewCachedCRUD() (*CachedCRUD, error) {
	var e error
	once.Do(
		func() {
			crud = new(CachedCRUD)
			redis, e := myredis.NewRedisConn()
			if e != nil {
				panic(e)
			}
			crud.redis = redis
			mysql, e := mysql.NewMysqlConn()
			if e != nil {
				panic(e)
			}
			crud.mysql = mysql
			oss, e := oss.GetOssService()
			if e != nil {
				panic(e)
			}
			crud.oss = oss
		},
	)

	return crud, e
}

func init() {
	var err error
	crud, err = NewCachedCRUD()
	if err != nil {
		panic(err)
	}
}

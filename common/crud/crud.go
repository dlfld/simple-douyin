package crud

import (
	"github.com/douyin/common/mysql"
	"github.com/douyin/common/oss/minioService"
	myredis "github.com/douyin/common/redis"
	"github.com/go-redis/redis/v8"

	"gorm.io/gorm"
)

type CachedCRUD struct {
	redis *redis.Client
	mysql *gorm.DB
	minio *minioService.MinioService
}

func NewCachedCRUD() (g *CachedCRUD, e error) {
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
	minio, e := minioService.GetMinio()
	if e != nil {
		return
	}
	g.minio = minio
	return
}

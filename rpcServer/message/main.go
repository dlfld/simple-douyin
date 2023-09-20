package main

import (
	"github.com/douyin/common/bloom"
	"github.com/douyin/common/etcd"
	"github.com/douyin/kitex_gen/model"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"log"
	"net"

	"github.com/cloudwego/kitex/server"
	"github.com/douyin/common/conf"
	"github.com/douyin/common/kafkaLog/productor"
	"github.com/douyin/common/mysql"
	rdb "github.com/douyin/common/redis"
	message "github.com/douyin/kitex_gen/message/messageservice"
	"gorm.io/gorm"
)

var db *gorm.DB
var cache *Cache
var bf *bloom.Filter

const (
	messageCacheTable  = "message"
	maxCacheMessageNum = 200
)

var logCollector *productor.LogCollector

func init() {
	var err error
	// 1. 初始化mysql
	if db, err = mysql.NewMysqlConn(); err != nil {
		panic(err)
	}
	if err = db.AutoMigrate(&model.Message{}); err != nil {
		panic(err)
	}
	// 2. 初始化redis缓存
	r, err := rdb.NewRedisConn()
	if err != nil {
		panic(err)
	}
	cache = &Cache{
		rdb: r,
	}
	// 3. 初始化日志收集器
	if logCollector, err = productor.NewLogCollector(conf.MessageService.Name); err != nil {
		panic(err)
	}
	//bf = bloom.NewBloom()
}

func main() {
	//p := otel.NewOtelProvider("message")
	//defer p.Shutdown(context.Background())
	addr, err := net.ResolveTCPAddr("tcp", conf.MessageService.Addr)
	if err != nil {
		log.Println(err.Error())
	}
	svr := message.NewServer(new(MessageServiceImpl), server.WithServiceAddr(addr), server.WithSuite(tracing.NewServerSuite()))
	etcd.RegisterService(conf.MessageService.Name, conf.MessageService.Addr)
	err = svr.Run()
	if err != nil {
		log.Println(err.Error())
	}

}

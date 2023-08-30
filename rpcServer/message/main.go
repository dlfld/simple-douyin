package main

import (
	"log"
	"net"

	"github.com/cloudwego/kitex/server"
	"github.com/douyin/common/conf"
	"github.com/douyin/common/jaeger"
	"github.com/douyin/common/etcd"
	"github.com/douyin/common/kafkaLog/productor"
	"github.com/douyin/common/mysql"
	rdb "github.com/douyin/common/redis"
	message "github.com/douyin/kitex_gen/message/messageservice"
	"github.com/douyin/kitex_gen/model"
	"gorm.io/gorm"
)

var db *gorm.DB
var cache *Cache

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
}

func main() {
	tracerSuite, closer := jaeger.InitJaegerServer("kitex-server-message")
	defer closer.Close()
	addr, err := net.ResolveTCPAddr("tcp", conf.MessageService.Addr)
	if err != nil {
		log.Println(err.Error())
	}
	svr := message.NewServer(new(MessageServiceImpl), server.WithServiceAddr(addr), server.WithSuite(tracerSuite))

	err = svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
	etcd.RegisterService(conf.MessageService.Name, conf.MessageService.Addr)
}

package main

import (
	"github.com/cloudwego/kitex/server"
	"github.com/douyin/common/comLogger"
	"github.com/douyin/common/conf"
	"github.com/douyin/common/mysql"
	rdb "github.com/douyin/common/redis"
	message "github.com/douyin/kitex_gen/message/messageservice"
	"github.com/douyin/kitex_gen/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"log"
	"net"
)

var db *gorm.DB
var cache *Cache
var err error
var logger *logrus.Logger

const messageCacheTable = "message"
const maxCacheMessageNum = 200

func init() {
	// 1. 初始化日志
	logger = comLogger.NewLogger()
	// 2. 初始化mysql
	db, err = mysql.NewMysqlConn()
	if err != nil {
		panic(err)
	}
	if err = db.AutoMigrate(&model.Message{}); err != nil {
		panic(err)
	}
	// 3. 初始化redis缓存
	r, err := rdb.NewRedisConn()
	if err != nil {
		panic(err)
	}
	cache = &Cache{
		rdb: r,
	}
}
func main() {
	addr, err := net.ResolveTCPAddr("tcp", conf.MessageService.Addr)
	svr := message.NewServer(new(MessageServiceImpl), server.WithServiceAddr(addr))

	err = svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}

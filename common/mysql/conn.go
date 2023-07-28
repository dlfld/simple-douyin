package mysql

import (
	"github.com/douyin/common/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

var db *gorm.DB
var once sync.Once

// NewMysqlConn 创建一个db数据库
func NewMysqlConn() (*gorm.DB, error) {
	var err error
	once.Do(func() {
		db, err = gorm.Open(mysql.Open(conf.Mysql.Login))
		if conf.Mysql.Debug {
			db = db.Debug()
		}
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

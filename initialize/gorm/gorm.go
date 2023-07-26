package initialize

import (
	"fmt"
	"os"

	"github.com/RaymondCode/simple-demo/models/config"
	"github.com/RaymondCode/simple-demo/models/db"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var _gorm *gorm.DB

/**
*  Gorm mysql数据库初始化
*  读取config/mysql.yaml中的数据 初始化mysql连接
 */
func GormInit() {
	// todo: 读取yaml文件的方式不优雅需要修改
	data, err := os.ReadFile("..\\..\\config\\mysql.yaml")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	conf := new(config.GeneralDB)
	//使用yaml.Unmarshal将yaml文件中的信息反序列化给Config结构体
	if err := yaml.Unmarshal(data, conf); err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	dsn := conf.Username + ":" + conf.Password + "@tcp(" + conf.Path + ":" + conf.Port + ")/" + conf.Dbname + "?" + conf.Config
	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}
	_gorm, err = gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   conf.Prefix,
			SingularTable: conf.Singular,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err == nil {
		_gorm.InstanceSet("gorm:table_options", "ENGINE="+conf.Engine)
		sqlDB, _ := _gorm.DB()
		sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
		sqlDB.SetMaxOpenConns(conf.MaxOpenConns)
	} else {
		// todo: 添加日志
		fmt.Println(err)
	}
}

// 创建数据表
func CreateTable() {
	if _gorm == nil {
		fmt.Println("err nil gorm")
		return
	}
	err := _gorm.AutoMigrate(
		db.Comment{},
		db.User{},
		db.Message{},
		db.Video{},
		db.FollowRelation{},
		db.FavoriteCommentRelation{},
		db.FavoriteVideoRelation{},
	)
	if err != nil {
		// todo: 添加日志
		fmt.Println(err)
	}
	// todo: 添加日志
	fmt.Println("create table success")
}

func Gorm() *gorm.DB {
	return _gorm
}

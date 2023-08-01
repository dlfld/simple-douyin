package initialize

import (
	"fmt"
	"log"

	"github.com/douyin/common/mysql"
	"github.com/douyin/models"
)

// 创建数据表
func CreateTable() {
	_gorm, err := mysql.NewMysqlConn()
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	err = _gorm.AutoMigrate(
		models.Comment{},
		models.User{},
		models.Message{},
		models.Video{},
		models.FollowRelation{},
		models.FavoriteCommentRelation{},
		models.FavoriteVideoRelation{},
	)
	if err != nil {
		// todo: 添加日志
		fmt.Println(err)
	}
	// todo: 添加日志
	fmt.Println("create table success")
}

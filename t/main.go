package main

import (
	"fmt"
	"github.com/douyin/common/mysql"
	"github.com/douyin/models"
)

func main() {
	conn, err := mysql.NewMysqlConn()
	if err != nil {
		return
	}
	conn.AutoMigrate(&models.FavoriteVideoRelation{})
	t := models.FavoriteVideoRelation{
		VideoID: 123,
		UserID:  222,
	}
	res := conn.Create(&t)
	fmt.Println(res.RowsAffected)
}

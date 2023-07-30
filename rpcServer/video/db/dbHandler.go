// Package db 用于数据库操作的方法
package db

import (
	mysql "github.com/douyin/common/mysql"
	"github.com/douyin/models"
)

// FindVideoListBy
// @Description: 根据输入的字段名和条件查询视频信息列表
// @param field: 字段名
// @param condition: 条件
// @return []models.Video: 视频信息列表
// @return error
func FindVideoListBy(field, condition string) ([]*models.Video, error) {
	conn, err := mysql.NewMysqlConn()
	if err != nil {
		return nil, err
	}
	videos := make([]*models.Video, 0)
	//根据field（表的字段）和指定的条件查询列表
	conn.Where(field+" = ?", condition).Find(&videos)
	//conn.
	return videos, nil
}

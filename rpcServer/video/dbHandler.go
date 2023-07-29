/*
*

	@author:戴林峰
	@date:2023/7/29
	@node:

*
*/
package main

import (
	"github.com/douyin/common/mysql"
	"github.com/douyin/kitex_gen/model"
	//"github.com/douyin/kitex_gen/model"
)

// FindVideoListBy
// @Description: 根据输入的字段名和条件查询视频信息列表
// @param field: 字段名
// @param condition: 条件
// @return []models.Video: 视频信息列表
// @return error
func FindVideoListBy(field, condition string) ([]*model.Video, error) {
	conn, err := mysql.NewMysqlConn()
	if err != nil {
		return nil, err
	}
	videos := make([]*model.Video, 0)
	//根据field（表的字段）和指定的条件查询列表
	conn.Where(field+" = ？", condition).Find(&videos)
	return videos, nil
}

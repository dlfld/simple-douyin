//
// Package db
// @Description: 数据库数据库操作业务逻辑
// @Author hehehhh
// @Date 2023-01-21 14:33:47
// @Update
//

package models

import (
	"github.com/douyin/common/mysql"
	"time"

	"gorm.io/gorm"
)

// Video
//
//	@Description: 视频数据模型
type Video struct {
	ID            int64     `gorm:"primarykey"`
	CreatedAt     time.Time `gorm:"not null;index:idx_create" json:"created_at,omitempty"`
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	Author        User           `gorm:"foreignkey:AuthorID" json:"author,omitempty"`
	AuthorID      int64          `gorm:"index:idx_authorid;not null" json:"author_id,omitempty"`
	PlayUrl       string         `gorm:"type:varchar(512);not null" json:"play_url,omitempty"`
	CoverUrl      string         `gorm:"type:varchar(255)" json:"cover_url,omitempty"`
	FavoriteCount int64          `gorm:"default:0;not null" json:"favorite_count,omitempty"`
	CommentCount  int64          `gorm:"default:0;not null" json:"comment_count,omitempty"`
	Title         string         `gorm:"type:varchar(50);not null" json:"title,omitempty"`
}

func (Video) TableName() string {
	return "videos"
}

// FindVideoListBy
// @Description: 根据输入的字段名和条件查询视频信息列表
// @param field: 字段名
// @param condition: 条件
// @return []models.Video: 视频信息列表
// @return error
func FindVideoListBy(field, condition string) ([]*Video, error) {
	conn, err := mysql.NewMysqlConn()
	if err != nil {
		return nil, err
	}
	videos := make([]*Video, 0)
	//根据field（表的字段）和指定的条件查询列表
	conn.Where(field+" = ?", condition).Find(&videos)
	//conn.
	return videos, nil
}

//	InsertVideo 插入数据
//
// @Description:
// @param video
// @return error
func InsertVideo(video *Video) error {
	conn, err := mysql.NewMysqlConn()
	if err != nil {
		return err
	}
	conn.Create(video)
	return nil
}

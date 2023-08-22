//
// Package db
// @Description: 数据库数据库操作业务逻辑
// @Author hehehhh
// @Date 2023-01-21 14:33:47
// @Update
//

package models

import (
	"time"
)

// Comment
//
//	@Description: 用户评论数据模型
type Comment struct {
	ID          int64     `gorm:"primarykey"`
	CreatedTime time.Time `json:"created_time"`
	//UpdatedAt time.Time
	//DeletedAt gorm.DeletedAt `gorm:"index"`
	//Video   Video `gorm:"foreignkey:VideoID" json:"video,omitempty"`
	VideoID int64 `gorm:"index:idx_videoid;not null" json:"video_id"`
	//User       User           `gorm:"foreignkey:UserID" json:"user,omitempty"`
	UserID  int64  `gorm:"index:idx_userid;not null" json:"user_id"`
	Content string `gorm:"type:varchar(255);not null" json:"content"`
}

func (Comment) TableName() string {
	return "comments"
}

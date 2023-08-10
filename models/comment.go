//
// Package db
// @Description: 数据库数据库操作业务逻辑
// @Author hehehhh
// @Date 2023-01-21 14:33:47
// @Update
//

package models

import "time"

// Comment
//
//	@Description: 用户评论数据模型
type Comment struct {
	ID int64 `gorm:"primarykey"`
	//CreatedAt time.Time `gorm:"index;not null" json:"create_date"`
	//UpdatedAt time.Time
	//DeletedAt gorm.DeletedAt `gorm:"index"`
	Video   Video `gorm:"foreignkey:VideoID" json:"video,omitempty"`
	VideoID int64 `gorm:"index:idx_videoid;not null" json:"video_id"`
	//User       User           `gorm:"foreignkey:UserID" json:"user,omitempty"`
	UserID     int64     `gorm:"index:idx_userid;not null" json:"user_id"`
	Content    string    `gorm:"type:varchar(255);not null" json:"content"`
	CreateTime time.Time `gorm:"type:datetime;not null" json:"create_time"`
	LikeCount  int64     `gorm:"column:like_count;default:0;not null" json:"like_count,omitempty"`
	TeaseCount int64     `gorm:"column:tease_count;default:0;not null" json:"tease_count,omitempty"`
}

func (Comment) TableName() string {
	return "comments"
}

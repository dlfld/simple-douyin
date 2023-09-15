//
// Package db
// @Description: 数据库数据库操作业务逻辑
// @Author hehehhh
// @Date 2023-01-21 14:33:47
// @Update
//

package models

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/douyin/common/gorse"
	"gorm.io/gorm"
)

// Comment
//
//	@Description: 用户评论数据模型
type Comment struct {
	ID          int64     `gorm:"primarykey"`
	CreatedTime time.Time `json:"created_time"`
	//UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	//Video   Video `gorm:"foreignkey:VideoID" json:"video,omitempty"`
	VideoID int64 `gorm:"index:idx_videoid;not null" json:"video_id"`
	//User       User           `gorm:"foreignkey:UserID" json:"user,omitempty"`
	UserID  int64  `gorm:"index:idx_userid;not null" json:"user_id"`
	Content string `gorm:"type:varchar(255);not null" json:"content"`
}

func (Comment) TableName() string {
	return "comments"
}

func (c *Comment) AfterUpdate(tx *gorm.DB) (err error) {
	cache.HDel(context.Background(), "UserInfoCache", fmt.Sprintf("%d", c.UserID))
	cache.Del(context.Background(), fmt.Sprintf("video:cache:%d", c.VideoID))
	return nil
}

func (c *Comment) AfterCreate(tx *gorm.DB) (err error) {
	gorse.Client.InsertFeedback(context.Background(),
		[]gorse.Feedback{{
			FeedbackType: "comment",
			UserId:       strconv.Itoa(int(c.UserID)),
			ItemId:       strconv.Itoa(int(c.VideoID)),
			Timestamp:    time.Now().Format("2006-01-02 15:04:05")}},
	)
	return nil
}

func (c *Comment) AfterDelete(tx *gorm.DB) (err error) {
	gorse.Client.DelFeedback(context.Background(), "comment",
		strconv.Itoa(int(c.UserID)), strconv.Itoa(int(c.VideoID)),
	)
	return nil
}

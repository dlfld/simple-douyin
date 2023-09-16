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

// FavoriteVideoRelation
//
//	@Description: 用户与视频的点赞关系数据模型
type FavoriteVideoRelation struct {
	VideoID int64 `gorm:"not null" json:"video_id"` //`gorm:"index:idx_videoid;not null" json:"video_id"`
	UserID  int64 `gorm:"not null" json:"user_id"`  //`gorm:"index:idx_userid;not null" json:"user_id"`
}

// FavoriteCommentRelation
//
//	@Description: 用户与评论的点赞关系数据模型
type FavoriteCommentRelation struct {
	//Comment   Comment `gorm:"foreignkey:CommentID;" json:"comment,omitempty"`
	CommentID uint `gorm:"column:comment_id;index:idx_commentid;not null" json:"video_id"`
	//User      User    `gorm:"foreignkey:UserID;" json:"user,omitempty"`
	UserID uint `gorm:"column:user_id;index:idx_userid;not null" json:"user_id"`
}

func (f *FavoriteVideoRelation) AfterUpdate(tx *gorm.DB) (err error) {
	cache.HDel(context.Background(), "UserInfoCache", fmt.Sprintf("%d", f.UserID))

	return cache.Del(context.Background(), fmt.Sprintf("video:cache:%d", f.VideoID)).Err()
}

func (f *FavoriteVideoRelation) AfterCreate(tx *gorm.DB) (err error) {
	gorse.Client.InsertFeedback(context.Background(),
		[]gorse.Feedback{{FeedbackType: "star",
			UserId:    strconv.Itoa(int(f.UserID)),
			ItemId:    strconv.Itoa(int(f.VideoID)),
			Timestamp: time.Now().Format("2006-01-02 15:04:05")}},
	)
	return nil
}
func (f *FavoriteVideoRelation) AfterDelete(tx *gorm.DB) (err error) {
	gorse.Client.DelFeedback(context.Background(),
		"star", strconv.Itoa(int(f.UserID)), strconv.Itoa(int(f.VideoID)),
	)
	return nil
}

func (FavoriteVideoRelation) TableName() string {
	return "user_favorite_videos"
}

func (FavoriteCommentRelation) TableName() string {
	return "user_favorite_comments"
}

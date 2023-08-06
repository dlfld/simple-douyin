//
// Package db
// @Description: 数据库数据库操作业务逻辑
// @Author hehehhh
// @Date 2023-01-21 14:33:47
// @Update
//

package models

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
	Comment   Comment `gorm:"foreignkey:CommentID;" json:"comment,omitempty"`
	CommentID uint    `gorm:"column:comment_id;index:idx_commentid;not null" json:"video_id"`
	//User      User    `gorm:"foreignkey:UserID;" json:"user,omitempty"`
	UserID uint `gorm:"column:user_id;index:idx_userid;not null" json:"user_id"`
}

func (FavoriteVideoRelation) TableName() string {
	return "user_favorite_videos"
}

func (FavoriteCommentRelation) TableName() string {
	return "user_favorite_comments"
}

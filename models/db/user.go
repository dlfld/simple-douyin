//
// Package db
// @Description: 数据库数据库操作业务逻辑
// @Author hehehhh
// @Date 2023-01-21 14:33:47
// @Update
//

package db

import (
	"gorm.io/gorm"
)

// User
//
//	@Description: 用户数据模型
type User struct {
	gorm.Model
	UserName        string  `gorm:"index:idx_username,unique;type:varchar(40);not null" json:"name,omitempty"`
	Password        string  `gorm:"type:varchar(256);not null" json:"password,omitempty"`
	FavoriteVideos  []Video `gorm:"many2many:user_favorite_videos" json:"favorite_videos,omitempty"`
	FollowingCount  uint    `gorm:"default:0;not null" json:"follow_count,omitempty"`                                                           // 关注总数
	FollowerCount   uint    `gorm:"default:0;not null" json:"follower_count,omitempty"`                                                         // 粉丝总数
	Avatar          string  `gorm:"type:varchar(256)" json:"avatar,omitempty"`                                                                  // 用户头像
	BackgroundImage string  `gorm:"column:background_image;type:varchar(256);default:default_background.jpg" json:"background_image,omitempty"` // 用户个人页顶部大图
	WorkCount       uint    `gorm:"default:0;not null" json:"work_count,omitempty"`                                                             // 作品数
	FavoriteCount   uint    `gorm:"default:0;not null" json:"favorite_count,omitempty"`                                                         // 喜欢数
	TotalFavorited  uint    `gorm:"default:0;not null" json:"total_favorited,omitempty"`                                                        // 获赞总量
	Signature       string  `gorm:"type:varchar(256)" json:"signature,omitempty"`                                                               // 个人简介
}

func (User) TableName() string {
	return "users"
}

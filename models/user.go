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
	"github.com/douyin/kitex_gen/model"
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

// CreateUser create user
func CreateUser(username, encryptPassword string) (*User, error) {
	newUser := User{UserName: username, Password: encryptPassword}
	db, err := mysql.NewMysqlConn()
	if err != nil {
		return nil, err
	}
	if result := db.Create(&newUser); result.Error != nil {
		return nil, result.Error
	}
	return &newUser, nil
}

// GetUserByName query user by name
func GetUserByName(username string) (*User, error) {
	db, err := mysql.NewMysqlConn()
	if err != nil {
		return nil, err
	}
	var user User
	if result := db.Where("user_name = ?", username).First(&user); result.Error != nil || user.ID <= 0 {
		return nil, result.Error
	}
	return &user, nil
}

// GetUserByUserId query user by name
func GetUserByUserId(userId uint) (*User, error) {
	db, err := mysql.NewMysqlConn()
	if err != nil {
		return nil, err
	}
	var user User
	if result := db.First(&user, userId); result.Error != nil || user.ID <= 0 {
		return nil, result.Error
	}
	return &user, nil
}

// ChangeModel model Change IsFollow UnSure!!!
func ChangeModel(u User) *model.User {
	FollowCount := int64(u.FollowingCount)
	FollowerCount := int64(u.FollowerCount)
	//IsFollow:=u
	Avatar := u.Avatar
	BackgroundImage := u.BackgroundImage
	Signature := u.Signature
	TotalFavorited := int64(u.TotalFavorited)
	WorkCount := int64(u.WorkCount)
	FavoriteCount := int64(u.FavoriteCount)

	user := model.User{
		Id:              int64(u.ID),
		Name:            u.UserName,
		FollowCount:     &FollowCount,
		FollowerCount:   &FollowerCount,
		IsFollow:        false,
		Avatar:          &Avatar,
		BackgroundImage: &BackgroundImage,
		Signature:       &Signature,
		TotalFavorited:  &TotalFavorited,
		WorkCount:       &WorkCount,
		FavoriteCount:   &FavoriteCount,
	}
	return &user
}

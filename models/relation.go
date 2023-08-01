//
// Package db
// @Description: 数据库数据库操作业务逻辑
// @Author hehehhh
// @Date 2023-01-21 14:33:47
// @Update
//

package models

import (
	"errors"

	"github.com/douyin/common/mysql"
	"gorm.io/gorm"
)

// FollowRelation
//
//	@Description: 用户之间的关注关系数据模型
type FollowRelation struct {
	gorm.Model
	User     User `gorm:"foreignkey:UserID;" json:"user,omitempty"`
	UserID   uint `gorm:"index:idx_userid;not null" json:"user_id"`
	ToUser   User `gorm:"foreignkey:ToUserID;" json:"to_user,omitempty"`
	ToUserID uint `gorm:"index:idx_userid;index:idx_userid_to;not null" json:"to_user_id"`
}

func (FollowRelation) TableName() string {
	return "relations"
}

func Follow(userID, toUserID uint) (err error) {
	db, _ := mysql.NewMysqlConn()

	var user User

	relation := FollowRelation{
		UserID:   userID,
		ToUserID: toUserID,
	}
	var n int64
	db.Find(&user, toUserID).Count(&n)
	if n == 0 {
		err = errors.New("user not found")
		return
	}
	err = db.Where("to_user_id=? and user_id = ?", toUserID, userID).FirstOrCreate(&relation).Error

	return
}

func Unfollow(userID, toUserID uint) (err error) {
	db, _ := mysql.NewMysqlConn()

	var toUser User

	relation := FollowRelation{
		UserID:   userID,
		ToUserID: toUserID,
	}
	var n int64
	db.Find(&toUser, toUserID).Count(&n)
	if n == 0 {
		err = errors.New("user not found")
		return
	}

	err = db.Where("to_user_id=? and user_id = ?", toUserID, userID).Delete(&relation).Error
	return
}

func GetFollowList(userID uint) (userList []*User, err error) {
	var relations []*FollowRelation
	db, _ := mysql.NewMysqlConn()

	err = db.Where("user_id=?", userID).Preload("ToUser").Find(&relations).Error
	if err != nil {
		return
	}
	for _, v := range relations {
		userList = append(userList, &v.ToUser)
	}
	return
}

func GetFollowerList(userID uint) (userList []*User, err error) {
	var relations []*FollowRelation
	db, _ := mysql.NewMysqlConn()

	err = db.Where("to_user_id=?", userID).Preload("User").Find(&relations).Error
	if err != nil {
		return
	}
	for _, v := range relations {
		userList = append(userList, &v.User)
	}
	return
}

func GetFriendList(userID uint) (userList []*User, err error) {
	var relations []*FollowRelation
	db, _ := mysql.NewMysqlConn()
	// redis, _ := redis.NewRedisConn()

	// select * from relations where user_id=1 and to_user_id in (select user_id from relations where to_user_id=1);
	subQuery := db.Model(&FollowRelation{}).Select("user_id").Where("to_user_id = ?", userID)
	db.Where("user_id = ? AND to_user_id IN (?)", userID, subQuery).Preload("ToUser").Find(&relations)
	if err != nil {
		return
	}
	for _, v := range relations {
		userList = append(userList, &v.ToUser)
	}
	return
}

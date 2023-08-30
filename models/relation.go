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

var db *gorm.DB

func init() {
	var err error
	db, err = mysql.NewMysqlConn()
	if err != nil {
		panic(err)
	}

}

// FollowRelation
//
//	@Description: 用户之间的关注关系数据模型
type FollowRelation struct {
	User     User `gorm:"foreignkey:UserID;" json:"user,omitempty"`
	UserID   uint `gorm:"index:idx_userid;not null" json:"user_id"`
	ToUser   User `gorm:"foreignkey:ToUserID;" json:"to_user,omitempty"`
	ToUserID uint `gorm:"index:idx_userid;index:idx_userid_to;not null" json:"to_user_id"`
}

func (FollowRelation) TableName() string {
	return "relations"
}

func Follow(userID, toUserID uint) (err error) {

	var user, toUser User

	relation := FollowRelation{
		UserID:   userID,
		ToUserID: toUserID,
	}

	db.Take(&user, userID)
	db.Take(&toUser, toUserID)
	if user.ID != userID || toUser.ID != toUserID {
		err = errors.New("user not found")
		return
	}
	err = db.Transaction(func(tx *gorm.DB) (err error) {
		// 更新关注数
		err = tx.Model(&user).Update("following_count", gorm.Expr("following_count + ?", 1)).Error
		if err != nil {
			return nil
		}
		err = tx.Model(&toUser).Update("follower_count", gorm.Expr("follower_count + ?", 1)).Error
		if err != nil {
			return nil
		}
		err = tx.Where("to_user_id=? and user_id = ?", toUserID, userID).FirstOrCreate(&relation).Error
		if err != nil {
			return nil
		}
		return nil
	})
	return
}

// Unfollow
//
//	@Description: 用户取消关注, 更新Mysql与Redis缓存
func Unfollow(userID, toUserID uint) (err error) {
	// cache,_:=redis.NewRedisConn()
	var toUser, user User

	relation := FollowRelation{
		UserID:   userID,
		ToUserID: toUserID,
	}
	db.Take(&user, userID)
	db.Take(&toUser, toUserID)
	if user.ID != userID || toUser.ID != toUserID {
		err = errors.New("user not found")
		return
	}
	db.Transaction(func(tx *gorm.DB) (err error) {
		// 更新关注数
		err = tx.Model(&user).Update("following_count", gorm.Expr("following_count - ?", 1)).Error
		if err != nil {
			return nil
		}
		err = tx.Model(&toUser).Update("follower_count", gorm.Expr("follower_count - ?", 1)).Error
		if err != nil {
			return nil
		}
		err = tx.Where("to_user_id=? and user_id = ?", toUserID, userID).Delete(&relation).Error
		if err != nil {
			return nil
		}
		return nil
	})

	return
}

func GetFollowList(userID uint) (userList []*User, err error) {
	if userID <= 0 {
		err = errors.New("user id is null")
		return
	}
	var relations []*FollowRelation

	err = db.Where("user_id=?", userID).Preload("ToUser").Find(&relations).Error
	if err != nil {
		return
	}
	userList = make([]*User, len(relations))
	for i, v := range relations {
		userList[i] = &v.ToUser
	}
	return
}

func GetFollowerList(userID uint) (userList []*User, err error) {
	if userID <= 0 {
		err = errors.New("user id is null")
		return
	}
	var relations []*FollowRelation

	err = db.Where("to_user_id=?", userID).Preload("User").Find(&relations).Error
	if err != nil {
		return
	}
	userList = make([]*User, len(relations))
	for i, v := range relations {
		userList[i] = &v.User
	}
	return
}

func GetFriendList(userID uint) (userList []*User, err error) {
	if userID <= 0 {
		err = errors.New("user id is null")
		return
	}
	var relations []*FollowRelation

	// redis, _ := redis.NewRedisConn()

	// select * from relations where user_id=1 and to_user_id in (select user_id from relations where to_user_id=1);
	subQuery := db.Model(&FollowRelation{}).Select("user_id").Where("to_user_id = ?", userID)
	db.Where("user_id = ? AND to_user_id IN (?)", userID, subQuery).Preload("ToUser").Find(&relations)
	if err != nil {
		return
	}
	userList = make([]*User, len(relations))
	for i, v := range relations {
		userList[i] = &v.ToUser
	}
	return
}

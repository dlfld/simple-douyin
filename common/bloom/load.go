package bloom

import (
	"gorm.io/gorm"
)

// 预先加载

func (b *Filter) PreLoadAll(db *gorm.DB) {
	b.PreLoadVideoIds(db)
	b.PreLoadCommentIds(db)
	b.PreLoadUserIds(db)
}

func (b *Filter) PreLoadVideoIds(db *gorm.DB) {
	result := db.Raw("SELECT id from videos")
	var videoIds []int64
	result.Scan(&videoIds)
	b.AddVideoIds(videoIds)
}

func (b *Filter) PreLoadCommentIds(db *gorm.DB) {
	result := db.Raw("SELECT id from comments")
	var commentIds []int64
	result.Scan(&commentIds)
	b.AddVideoIds(commentIds)
}

func (b *Filter) PreLoadUserIds(db *gorm.DB) {
	result := db.Raw("SELECT id from users")
	var userIds []int64
	result.Scan(&userIds)
	b.AddVideoIds(userIds)
}

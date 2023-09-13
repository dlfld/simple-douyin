package bloom

import (
	"gorm.io/gorm"
	"strconv"
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
	b.AddElements(convertInt64ToPrefixString(videoIds, videoPrefix))
}

func (b *Filter) PreLoadCommentIds(db *gorm.DB) {
	result := db.Raw("SELECT id from comments")
	var commentIds []int64
	result.Scan(&commentIds)
	b.AddElements(convertInt64ToPrefixString(commentIds, commentPrefix))
}

func (b *Filter) PreLoadUserIds(db *gorm.DB) {
	result := db.Raw("SELECT id from users")
	var userIds []int64
	result.Scan(&userIds)
	b.AddElements(convertInt64ToPrefixString(userIds, userPrefix))
}

func convertInt64ToPrefixString(intSlice []int64, prefix string) []string {
	stringSlice := make([]string, len(intSlice))
	for i, num := range intSlice {
		stringSlice[i] = prefix + strconv.FormatInt(num, 10)
	}
	return stringSlice
}

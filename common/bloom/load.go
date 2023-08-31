package bloom

import (
	"gorm.io/gorm"
	"strconv"
)

func (b *Filter) PreLoadAll(db *gorm.DB) {
	go b.PreLoadVideoIds(db)
	go b.PreLoadCommentIds(db)
	go b.PreLoadUserIds(db)
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

func (b *Filter) AddElements(elements []string) {
	for _, e := range elements {
		b.filter.Add([]byte(e))
	}
}

func convertInt64ToPrefixString(intSlice []int64, prefix string) []string {
	stringSlice := make([]string, len(intSlice))
	for i, num := range intSlice {
		stringSlice[i] = prefix + strconv.FormatInt(num, 10)
	}
	return stringSlice
}

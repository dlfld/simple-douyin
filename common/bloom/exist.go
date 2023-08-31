package bloom

import "strconv"

// IfExists 通用查询, element需自定义, 慎用
func (b *Filter) IfExists(element string) bool {
	return b.filter.Test([]byte(element))
}

func (b *Filter) IfVideoIdExists(videoId int64) bool {
	element := videoPrefix + strconv.FormatInt(videoId, 10)
	return b.filter.Test([]byte(element))
}

func (b *Filter) IfCommentIdExists(commentId int64) bool {
	element := commentPrefix + strconv.FormatInt(commentId, 10)
	return b.filter.Test([]byte(element))
}

func (b *Filter) IfUserIdExists(userId int64) bool {
	element := userPrefix + strconv.FormatInt(userId, 10)
	return b.filter.Test([]byte(element))
}

package bloom

// IfExists 通用查询, element需自定义, 慎用
func (b *Filter) IfExists(element string) bool {
	return b.filter.Test([]byte(element))
}

func (b *Filter) IfVideoIdExists(videoId int64) bool {
	element := getVideoElement(videoId)
	return b.filter.Test([]byte(element))
}

func (b *Filter) IfCommentIdExists(commentId int64) bool {
	element := getCommentElement(commentId)
	return b.filter.Test([]byte(element))
}

func (b *Filter) IfUserIdExists(userId int64) bool {
	element := getUserElement(userId)
	return b.filter.Test([]byte(element))
}

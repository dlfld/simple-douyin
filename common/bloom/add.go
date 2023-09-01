package bloom

func (b *Filter) AddElements(elements []string) {
	for _, e := range elements {
		b.filter.Add([]byte(e))
	}
}

func (b *Filter) AddVideoIds(videoId int64) {
	element := getVideoElement(videoId)
	b.filter.Add([]byte(element))
}

func (b *Filter) AddCommentIds(commentId int64) {
	element := getCommentElement(commentId)
	b.filter.Add([]byte(element))
}

func (b *Filter) AddUserIds(userId int64) {
	element := getUserElement(userId)
	b.filter.Add([]byte(element))
}

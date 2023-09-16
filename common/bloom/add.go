package bloom

import "context"

const (
	bloomVideoId   = "bloom:videoIds"
	bloomCommentId = "bloom:commentIds"
	bloomUserId    = "bloom:userIds"
	bloomUserName  = "bloom:userNames"
)

func (b *Filter) addElement(key string, element interface{}) {
	b.filter.Do(context.Background(), "BF.ADD", key, element)
}

func (b *Filter) AddVideoId(videoId int64) {
	b.addElement(bloomVideoId, videoId)
}

func (b *Filter) AddCommentId(commentId int64) {
	b.addElement(bloomCommentId, commentId)
}

func (b *Filter) AddUserId(userId int64) {
	b.addElement(bloomUserId, userId)
}

func (b *Filter) AddUserName(name string) {
	b.addElement(bloomUserName, name)
}

func (b *Filter) addElements(key string, elements []interface{}) {
	for _, e := range elements {
		b.addElement(key, e)
	}
}

func (b *Filter) AddVideoIds(elements []int64) {
	for _, e := range elements {
		b.addElement(bloomVideoId, e)
	}
}

func (b *Filter) AddCommentIds(elements []int64) {
	for _, e := range elements {
		b.addElement(bloomCommentId, e)
	}
}

func (b *Filter) AddUserIds(elements []int64) {
	for _, e := range elements {
		b.addElement(bloomUserId, e)
	}
}

func (b *Filter) AddUserNames(elements []string) {
	for _, e := range elements {
		b.addElement(bloomUserName, e)
	}
}

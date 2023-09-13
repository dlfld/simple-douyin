package bloom

import "context"

const (
	bloomVideo   = "bloom:videoIds"
	bloomComment = "bloom:commentIds"
	bloomUser    = "bloom:userIds"
)

func (b *Filter) addElement(key string, element interface{}) {
	b.filter.Do(context.Background(), "BF.ADD", key, element)
}

func (b *Filter) AddVideoId(videoId int64) {
	b.addElement(bloomVideo, videoId)
}

func (b *Filter) AddCommentId(commentId int64) {
	b.addElement(bloomComment, commentId)
}

func (b *Filter) AddUserId(userId int64) {
	b.addElement(bloomUser, userId)
}

func (b *Filter) addElements(key string, elements []interface{}) {
	for _, e := range elements {
		b.addElement(key, e)
	}
}

func (b *Filter) AddVideoIds(elements []int64) {
	for _, e := range elements {
		b.addElement(bloomVideo, e)
	}
}

func (b *Filter) AddCommentIds(elements []int64) {
	for _, e := range elements {
		b.addElement(bloomComment, e)
	}
}

func (b *Filter) AddUserIds(elements []int64) {
	for _, e := range elements {
		b.addElement(bloomUser, e)
	}
}

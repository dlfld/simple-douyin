package bloom

import (
	"context"
)

func (b *Filter) ifExists(key string, element interface{}) (exist bool, err error) {
	result, err := b.filter.Do(context.Background(), "BF.EXISTS", key, element).Bool()
	if err != nil {
		return
	} else {
		return result, nil
	}
}

func (b *Filter) CheckIfVideoIdExists(videoId int64) (exist bool, err error) {
	return b.ifExists(bloomVideoId, videoId)
}

func (b *Filter) CheckIfCommentIdExists(commentId int64) (exist bool, err error) {
	return b.ifExists(bloomCommentId, commentId)
}

func (b *Filter) CheckIfUserIdExists(userId int64) (exist bool, err error) {
	return b.ifExists(bloomUserId, userId)
}

func (b *Filter) CheckIfUserNameExists(userName string) (exist bool, err error) {
	return b.ifExists(bloomUserName, userName)
}

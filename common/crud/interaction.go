package crud

import "log"

// IsFavorite 判断是否点赞
func IsFavorite(self uint, videoId uint) (isFavorite bool, err error) {
	var i int64
	result := crud.mysql.Raw("select 1 from user_favorite_videos WHERE user_id = ? AND video_id = ? LIMIT 1", self, videoId).Count(&i)

	if result.Error != nil {
		log.Println(err)
	}
	return i != 0, result.Error
}

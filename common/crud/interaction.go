package crud

import (
	"fmt"
	_ "github.com/douyin/kitex_gen/interaction"
	"github.com/douyin/kitex_gen/model"
	"github.com/douyin/models"
)

func (c *CachedCRUD) SearchFavoriteExist(m *models.FavoriteVideoRelation) (exist bool, err error) {
	result := c.mysql.Where("user_id = ? and video_id = ?", m.UserID, m.VideoID).First(m)
	if result.Error != nil {
		// TODO : log err
		fmt.Println(result.Error)
	}
	return result.RowsAffected > 0, result.Error
}

func (c *CachedCRUD) InsertFavorite(m *models.FavoriteVideoRelation) (rows int64, err error) {

	result := c.mysql.Create(m)
	if result.Error != nil {
		// TODO : log err
		fmt.Println(result.Error)
	}
	return result.RowsAffected, result.Error

}

func (c *CachedCRUD) CancelFavorite(m *models.FavoriteVideoRelation) (rows int64, err error) {

	result := c.mysql.Where("user_id = ? and video_id = ?", m.UserID, m.VideoID).Delete(m)
	return result.RowsAffected, result.Error

}
func (c *CachedCRUD) SearchVideoListById(id int64) (videoList []*models.Video, err error) {
	result := c.mysql.Raw("SELECT * FROM videos WHERE id in (SELECT video_id from user_favorite_videos WHERE user_id = ?)", id)
	var t []*models.Video
	result.Scan(&t)
	return t, nil
}

func (c *CachedCRUD) SearchUserById(id int64) (user *model.User, err error) {
	result := c.mysql.Raw("SELECT * FROM users WHERE id = ?", id)
	var t model.User
	result.Scan(&t)
	return &t, err
}

// 评论

func (c *CachedCRUD) InsertComment(m *models.Comment) (rows int64, err error) {
	result := c.mysql.Create(m)
	if result.Error != nil {
		// TODO : log err
		fmt.Println(result.Error)
	}
	return result.RowsAffected, result.Error
}

func (c *CachedCRUD) DeleteComment(m *models.Comment) (rows int64, err error) {
	result := c.mysql.Where("id = ? and video_id = ?", m.ID, m.VideoID).Delete(m)
	return result.RowsAffected, result.Error
}

func (c *CachedCRUD) SearchCommentListSort(videoId int64) (videoList []*models.Comment, err error) {
	//result := c.mysql.Raw("SELECT * from comments c join users u on c.user_id = u.id  WHERE video_id = 3 ORDER BY c.id DESC")
	result := c.mysql.Raw("SELECT * from comments where video_id = ? ORDER BY id DESC", videoId)
	var t []*models.Comment
	result.Scan(&t)
	return t, nil
}

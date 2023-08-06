package crud

import (
	"fmt"
	_ "github.com/douyin/kitex_gen/interaction"
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

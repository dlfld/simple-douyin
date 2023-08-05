package crud

import (
	_ "github.com/douyin/kitex_gen/interaction"
	"github.com/douyin/models"
)

func (c *CachedCRUD) InsertFavorite(m *models.FavoriteVideoRelation) (rows int64, err error) {

	result := c.mysql.Create(m)
	return result.RowsAffected, result.Error

}

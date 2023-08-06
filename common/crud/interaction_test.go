package crud

import (
	"fmt"
	"github.com/douyin/models"
	"testing"
)

func TestInsertFavorite(t *testing.T) {
	c, err := NewCachedCRUD()
	if err != nil {
		return
	}
	m := models.FavoriteVideoRelation{
		VideoID: 100,
		UserID:  34324324,
	}
	rows, err := c.InsertFavorite(&m)
	if err != nil {
		return
	}
	fmt.Println("rows: ", rows)
}

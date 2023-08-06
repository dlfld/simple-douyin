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

func TestSearchVideoListById(t *testing.T) {
	c, err := NewCachedCRUD()
	if err != nil {
		return
	}
	videoList, err := c.SearchVideoListById(1)
	if err != nil {
		return
	}
	fmt.Println("videoList: ", videoList)
}

func TestSearchSearchUserById(t *testing.T) {
	c, err := NewCachedCRUD()
	if err != nil {
		return
	}
	user, err := c.SearchUserById(1)
	if err != nil {
		return
	}
	fmt.Println("user: ", user)
}

func TestDeleteComment(t *testing.T) {
	c, err := NewCachedCRUD()
	if err != nil {
		return
	}
	m := models.Comment{
		ID:      9,
		VideoID: 3,
	}
	rows, err := c.DeleteComment(&m)
	if err != nil {
		return
	}
	fmt.Println("rows: ", rows)
}

func TestSearchCommentListSort(t *testing.T) {
	c, err := NewCachedCRUD()
	if err != nil {
		return
	}
	commentList, err := c.SearchCommentListSort(3)
	if err != nil {
		return
	}
	fmt.Println("commentList: ", commentList)
	fmt.Println("time: ", commentList[0].CreateTime.GoString())
}

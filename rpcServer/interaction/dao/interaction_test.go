package dao

import (
	"fmt"
	"github.com/douyin/models"
	"testing"
)

func TestCreate(t *testing.T) {
	c := NewMysql()
	_ = c.cli.AutoMigrate(
		models.Comment{},
	)
}

func TestInsertFavorite(t *testing.T) {
	c := NewMysql()
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
	c := NewMysql()
	videoList, err := c.SearchVideoListById(1)
	if err != nil {
		return
	}
	fmt.Println("videoList: ", videoList)
}

func TestSearchUserById(t *testing.T) {
	c := NewMysql()
	user, err := c.SearchUserByAuthorId(1, 2)
	if err != nil {
		return
	}
	fmt.Println("user: ", user)
}

func TestDeleteComment(t *testing.T) {
	c := NewMysql()
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
	c := NewMysql()
	commentList, err := c.SearchCommentListSort(3)
	if err != nil {
		return
	}
	fmt.Println("commentList: ", commentList)
	fmt.Println("time: ", commentList[0].CreatedTime)
}

func TestSearchUserByAuthorIds(t *testing.T) {
	c := NewMysql()
	authorIds := []int64{1, 2, 3}
	userID := 2
	res, err := c.SearchUserByAuthorIds(authorIds, int64(userID))
	if err != nil {
		return
	}
	fmt.Println("res: ", res)
}

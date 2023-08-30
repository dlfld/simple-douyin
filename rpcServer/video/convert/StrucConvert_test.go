// @author:戴林峰
// @date:2023/7/30
// @node:
package convert

import (
	"fmt"
	"github.com/douyin/models"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestVideoSliceBo2Dto(t *testing.T) {
	bo := models.Video{
		ID:            1,
		CreatedAt:     time.Time{},
		UpdatedAt:     time.Time{},
		DeletedAt:     gorm.DeletedAt{},
		Author:        models.User{},
		AuthorID:      1,
		PlayUrl:       "http",
		CoverUrl:      "http",
		FavoriteCount: 2,
		CommentCount:  2,
		Title:         "测试title",
	}
	boSlice := []*models.Video{&bo}
	dto, _ := VideoSliceBo2Dto(boSlice)
	fmt.Printf("%v\n", dto[0])
}

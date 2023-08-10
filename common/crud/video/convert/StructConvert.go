// Package convert 用户结构体转换的方法
package convert

import (
	"github.com/douyin/kitex_gen/model"
	"github.com/douyin/models"
	"github.com/jinzhu/copier"
)

// VideoSliceBo2Dto
//
//		 @Description: 将视频bo转为视频dto列表
//		 @param boSlice bo列表
//	     @param dtoSlice dto列表
//		 @return 返回视频的dto列表，也就是kitex中的model
func VideoSliceBo2Dto(boSlice []*models.Video) ([]*model.Video, error) {
	dtoSlice := make([]*model.Video, 0, len(boSlice))
	// 将数据库对应的结构体转换为kitex对应的结构体
	for _, videoBo := range boSlice {
		videoDto := model.Video{}
		//对同名属性的转换，其中还有一个id是不同名的需要手动转换
		err := copier.Copy(&videoDto, &videoBo)
		// 两个结构体还有这个变量是不同名的
		videoDto.Id = videoBo.ID
		if err != nil {
			return nil, err
		}
		dtoSlice = append(dtoSlice, &videoDto)
	}
	return dtoSlice, nil
}

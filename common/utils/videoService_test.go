// @author:戴林峰
// @date:2023/8/5
// @node:
package utils

import (
	"fmt"
	"testing"
)

func TestFindVideoListByUserId(t *testing.T) {
	id, _ := FindVideoListByUserId(1)
	for _, item := range id {
		fmt.Printf("%v\n", item)
	}

}

func TestUploadVideo(t *testing.T) {
	//reader := bytes.NewReader([]byte{'1'})
	//// 上传文件的文件名
	//filename := "filename"
	//userId := 1
	//// TODO 根据Token获取用户信息，然后根据用户信息写入用户投稿的视频，在redis中加入这一条视频
	//videoUrl := fmt.Sprintf("http://%s/%s/%s", conf.MinioConfig.IP, conf.MinioConfig.VideoBucketName, filename)
	// TODO 魔法值需要改
	//contentType := "application/mp4"
	//title := "title"
	//dataLen := int64(len([]byte{'1'}))
	//UploadVideo(reader, filename, videoUrl, dataLen, int64(userId), title)
}

//func TestGetVideoFeed(t *testing.T) {
//	feed, _, _ := GetVideoFeed(0, 2)
//	for _, item := range feed {
//		fmt.Printf("%v\n", item)
//	}
//
//}

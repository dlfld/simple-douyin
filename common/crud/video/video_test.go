// @author:戴林峰
// @date:2023/8/5
// @node:
package video

import (
	"bytes"
	"fmt"
	"github.com/douyin/common/conf"
	"testing"
)

func TestFindVideoListByUserId(t *testing.T) {
	id, _ := FindVideoListByUserId(1)
	fmt.Printf("%+v", id)
}

func TestUploadVideo(t *testing.T) {
	reader := bytes.NewReader([]byte{'1'})
	// 上传文件的文件名
	filename := "filename"
	userId := 1
	// TODO 根据Token获取用户信息，然后根据用户信息写入用户投稿的视频，在redis中加入这一条视频
	videoUrl := fmt.Sprintf("http://%s/%s/%s", conf.MinioConfig.IP, conf.MinioConfig.VideoBucketName, filename)
	// TODO 魔法值需要改
	contentType := "application/mp4"
	title := "title"
	dataLen := int64(len([]byte{'1'}))
	UploadVideo(reader, filename, contentType, videoUrl, dataLen, int64(userId), title)
}

// @author:戴林峰
// @date:2023/8/2
// @node:

package main

import (
	"bytes"
	"github.com/douyin/common/conf"
	"github.com/douyin/common/oss"

	"io/ioutil"
	"testing"
	"time"
)

func TestVideoServiceImpl_PublishAction(t *testing.T) {
	videoBytes, _ := ioutil.ReadFile("vedio.mp4")

	service, _ := oss.GetOssService()
	reader := bytes.NewReader(videoBytes)
	// 上传文件的文件名
	filename := "vedio_" + time.Now().String() + ".mp4"
	// TODO 魔法值需要改
	contentType := "application/mp4"
	_ = service.UploadFileWithBytestream(conf.MinioConfig.VideoBucketName, reader, filename, int64(len(videoBytes)), contentType)

}

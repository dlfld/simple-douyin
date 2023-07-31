// @author:戴林峰
// @date:2023/7/30
// @node:
package main

import (
	"github.com/douyin/common/conf"
	"testing"
)

func TestGetClient(t *testing.T) {
	_, err := GetClient()
	if err != nil {
		return
	}
}

func TestCreateBucket(t *testing.T) {
	err := CreateBucket(conf.MinioConfig.VideoBucketName)
	if err != nil {
		return
	}
	err1 := CreateBucket(conf.MinioConfig.AvatarBucketName)
	if err1 != nil {
		return
	}
	err2 := CreateBucket(conf.MinioConfig.BackgroundImageBucketName)
	if err2 != nil {
		return
	}

}

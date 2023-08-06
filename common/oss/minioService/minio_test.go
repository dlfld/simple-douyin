// @author:戴林峰
// @date:2023/7/30
// @node:
package minioService

import (
	"github.com/douyin/common/conf"
	"testing"
)

func TestCreateBucket(t *testing.T) {
	minio, _ := GetMinio()
	err := minio.CreateBucket(conf.MinioConfig.VideoBucketName)
	if err != nil {
		return
	}
	err1 := minio.CreateBucket(conf.MinioConfig.AvatarBucketName)
	if err1 != nil {
		return
	}
	err2 := minio.CreateBucket(conf.MinioConfig.BackgroundImageBucketName)
	if err2 != nil {
		return
	}

}

func TestMinio_DeleteBucket(t *testing.T) {
	minio, _ := GetMinio()
	_ = minio.DeleteBucket(conf.MinioConfig.VideoBucketName)
	_ = minio.DeleteBucket(conf.MinioConfig.AvatarBucketName)
	_ = minio.DeleteBucket(conf.MinioConfig.BackgroundImageBucketName)
}

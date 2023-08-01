// @author:戴林峰
// @date:2023/8/1
// @node:
package oss

import (
	"fmt"
	"testing"
	"log"
)

func TestService_CreateBucket(t *testing.T) {
	ossService, _ := GetOssService()
	_ = ossService.CreateBucket("test1")
}
func TestGetOssService(t *testing.T) {
	service, _ := GetOssService()
	fmt.Printf("%+v", service)
}

func TestService_DeleteBucket(t *testing.T) {
	service, _ := GetOssService()
	_ = service.DeleteBucket("test1")
}

// func TestUploadVideo(t *testing.T) {
//     // 创建一个新的 Service 实例
// 	service, _ := GetOssService()

//     // 调用 UploadVideo 方法来上传一个视频文件
//     err := service.UploadVideo("test", "./minioService/data_test/bear.mp4")
//     // service.UploadVideo("test", "data_test/bear.mp4")

//     // 检查是否有任何错误
//     if err != nil {
//         t.Errorf("Failed to upload video: %v", err)
//     }
// }


// 上传头像图片
func TestUploadAvatar(t *testing.T) {
	// 创建一个新的 Service 实例
	service, err := GetOssService()
	if err != nil {
		t.Errorf("Failed to create OssService: %v", err)
	}

	err = service.UploadFile("test", "./minioService/data_test/头像1.jpeg", "image/jpeg")
	if err != nil {
		log.Printf("上传头像失败: %v\n", err)
	}
}

// 上传封面图片
func TestUploadCover(t *testing.T) {
    // 创建一个新的 Service 实例
	service, _ := GetOssService()

	var err error
	err = service.UploadFile("test", "./minioService/data_test/封面1.jpeg", "image/jpeg")
	if err != nil {
		log.Printf("上传封面失败: %v\n", err)
	}

}


// 上传视频
func TestUploadVideo(t *testing.T) {
    // 创建一个新的 Service 实例
	service, _ := GetOssService()

	err := service.UploadFile("test", "./minioService/data_test/bear.mp4", "video/mp4")
	if err != nil {
		log.Printf("上传视频文件失败: %v\n", err)
	}
}
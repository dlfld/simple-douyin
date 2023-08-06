// @author:戴林峰
// @date:2023/8/1
// @node:
package oss

import (
	"fmt"
	"testing"
	"log"
	"os"
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

	err = service.UploadFile("test", "./minioService/data_test/头像1.jpeg", "image/jpg")
	if err != nil {
		log.Printf("上传头像失败: %v\n", err)
	}
}


// 上传头像图片(传入字节流)
// 从用户上传的 HTTP 请求中获取 io.Reader 对象和文件的名字和大小
// file io.Reader
// fileInfo.Name() string
// fileInfo.Size() int64
func TestUploadAvatarWithBytestream(t *testing.T) {
    // 创建一个新的 Service 实例
	service, _ := GetOssService()

	// 打开图片文件
	file, err := os.Open("./minioService/data_test/头像1.jpg")
	if err != nil {
		log.Fatalf("打开图片失败: %v", err)
	}
	defer file.Close()

	// 获取文件大小
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("获取图片信息失败: %v", err)
	}

	// 上传图片
	err = service.UploadFileWithBytestream("test", file, fileInfo.Name(), fileInfo.Size(), "image/jpg")
	if err != nil {
		log.Fatalf("图片上传失败: %v", err)
	}

	log.Println("头像图片上传成功")

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

// 上传封面(传入字节流)
// 从用户上传的 HTTP 请求中获取 io.Reader 对象和文件的名字和大小
// file io.Reader
// fileInfo.Name() string
// fileInfo.Size() int64
func TestUploadCoverWithBytestream(t *testing.T) {
    // 创建一个新的 Service 实例
	service, _ := GetOssService()

	// 打开图片文件
	file, err := os.Open("./minioService/data_test/封面1.jpeg")
	if err != nil {
		log.Fatalf("打开图片失败: %v", err)
	}
	defer file.Close()

	// 获取文件大小
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("获取图片信息失败: %v", err)
	}


	// 上传图片
	err = service.UploadFileWithBytestream("test", file, fileInfo.Name(), fileInfo.Size(), "image/jpeg")
	if err != nil {
		log.Fatalf("图片上传失败: %v", err)
	}

	log.Println("图片上传成功")

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


// 上传头像图片(传入字节流)
// 从用户上传的 HTTP 请求中获取 io.Reader 对象和文件的名字和大小
// file io.Reader
// fileInfo.Name() string
// fileInfo.Size() int64
func TestUploadVideoWithBytestream(t *testing.T) {
    // 创建一个新的 Service 实例
	service, _ := GetOssService()

	// 打开图片文件
	file, err := os.Open("./minioService/data_test/bear.mp4")
	if err != nil {
		log.Fatalf("打开视频失败: %v", err)
	}
	defer file.Close()

	// 获取文件大小
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("获取视频信息失败: %v", err)
	}

	// 上传图片
	err = service.UploadFileWithBytestream("test", file, fileInfo.Name(), fileInfo.Size(), "video/mp4")
	if err != nil {
		log.Fatalf("视频上传失败: %v", err)
	}

	log.Println("视频上传成功")

}
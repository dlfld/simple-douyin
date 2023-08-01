package oss

import (
	"github.com/douyin/common/oss/minioService"
	"sync"
	"io"
)

// OssInterface Oss 服务接口
type OssInterface interface {
	// GetClient 获取连接客户端
	GetClient() (interface{}, error)
	// CreateBucket 创建桶
	CreateBucket(bucketName string) error
	// DeleteBucket 删除桶
	DeleteBucket(bucketName string) error
	// // 上传Video文件
	// UploadVideo(bucketName string, filePath string) error
	// 上传文件
	UploadFile(bucketName string, filePath string, contentType string) error
	// 上传文件(传入字节流)
	UploadFileWithBytestream(bucketName string, reader io.Reader, fileName string, fileSize int64, contentType string) error

	
}

// Service
type Service struct {
	ossService interface{}
}

// TODO 用这个结构体实现具体的方法？
// 通过反射获取当前包下具体的实现软件，然后拉进来？
//然后使用可配置的方式配置具体的对象存储软件？

func GetOssService() (*Service, error) {
	var once sync.Once
	var service *Service
	var err error
	once.Do(func() {
		var minioCase *minioService.MinioService
		//TODO 获取Minio实例 这里涉及到了调用具体的minio的方法，（后面看看能否改为反射的方式）
		minioCase, err = minioService.GetMinio()
		service = &Service{}
		service.ossService = minioCase

	})
	if err != nil {
		return nil, err
	}
	return service, nil
}

func (service *Service) GetClient() interface{} {
	//TODO 这里涉及到了调用具体的minio的方法，（后面看看能否改为反射的方式）
	return service.ossService.(*minioService.MinioService).Client
}

func (service *Service) CreateBucket(bucketName string) error {
	return service.ossService.(*minioService.MinioService).CreateBucket(bucketName)
}

func (service *Service) DeleteBucket(bucketName string) error {
	return service.ossService.(*minioService.MinioService).DeleteBucket(bucketName)
}

// func (service *Service) UploadVideo(bucketName string, filePath string) error {
// 	return service.ossService.(*minioService.MinioService).UploadVideo(bucketName, filePath)
// }


func (service *Service) UploadFile(bucketName string, filePath string, contentType string) error {
	return service.ossService.(*minioService.MinioService).UploadFile(bucketName, filePath, contentType)
}


// UploadFileWithBytestream
// @Description: 通用的文件上传函数(传入字节流, 信息从用户上传的 HTTP 请求中获取)
// @param bucketName 桶名字
// @param filePath 文件路径
// @param contentType 文件类型(image/jpeg video/mp4)
// @return error
func (service *Service) UploadFileWithBytestream(bucketName string, reader io.Reader, fileName string, fileSize int64, contentType string) error {
	return service.ossService.(*minioService.MinioService).UploadFileWithBytestream(bucketName, reader, fileName, fileSize, contentType)
}

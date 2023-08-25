package oss

import (
	"io"
	"sync"

	"github.com/douyin/common/oss/minioService"
)

// OssInterface Oss 服务接口
type OssInterface interface {
	// GetClient 获取连接客户端
	GetClient() (interface{}, error)
	// CreateBucket 创建桶
	CreateBucket(bucketName string) error
	// DeleteBucket 删除桶
	DeleteBucket(bucketName string) error
	// UploadFile  上传文件
	UploadFile(bucketName string, filePath string, contentType string) error
	// UploadFileWithBytestream 上传文件(传入字节流)
	UploadFileWithBytestream(bucketName string, reader io.Reader, fileName string, fileSize int64, contentType string) error
	// RemoveObject 删除对象
	RemoveObject(bulkName, objectName string) error
}

// Service OssService结构体
type Service struct {
	ossService interface{}
}

// TODO 用这个结构体实现具体的方法？
// 通过反射获取当前包下具体的实现软件，然后拉进来？
// 然后使用可配置的方式配置具体的对象存储软件？
var once sync.Once
var service *Service
var minioCase *minioService.MinioService

func init() {
	GetOssService()
}

//	GetOssService
//
// @Description: 获取OSS服务端
// @return *Service

// @return error
func GetOssService() (*Service, error) {
	var err error
	once.Do(func() {
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

//	GetClient
//
// @Description: 获取客户端
// @receiver service
// @return interface{}
// @return error
func (service *Service) GetClient() (interface{}, error) {
	//TODO 这里涉及到了调用具体的minio的方法，（后面看看能否改为反射的方式）
	return service.ossService.(*minioService.MinioService).Client, nil
}

//	CreateBucket
//
// @Description: 根据桶名创建桶
// @receiver service
// @param bucketName
// @return error
func (service *Service) CreateBucket(bucketName string) error {
	return service.ossService.(*minioService.MinioService).CreateBucket(bucketName)
}

//	DeleteBucket
//
// @Description: 根据桶名删除桶
// @receiver service
// @param bucketName
// @return error
func (service *Service) DeleteBucket(bucketName string) error {
	return service.ossService.(*minioService.MinioService).DeleteBucket(bucketName)
}

// UploadFile
// @Description: 上传文件
// @receiver service
// @param bucketName 桶名字
// @param filePath 文件位置
// @param contentType  文件类型(image/jpeg video/mp4)
// @return error
func (service *Service) UploadFile(bucketName string, filePath string, contentType string) error {
	return service.ossService.(*minioService.MinioService).UploadFile(bucketName, filePath, contentType)
}

// UploadFileWithBytestream
// @Description: 通用的文件上传函数(传入字节流, 信息从用户上传的 HTTP 请求中获取)
// @param bucketName 桶名字
// @param reader io.Reader
// @param fileName 文件名字
// @param fileSize 文件大小
// @param contentType 文件类型(image/jpeg video/mp4)
// @return error
func (service *Service) UploadFileWithBytestream(bucketName string, reader io.Reader, fileName string, fileSize int64, contentType string) error {
	return service.ossService.(*minioService.MinioService).UploadFileWithBytestream(bucketName, reader, fileName, fileSize, contentType)
}

// RemoveObject
// @Description: 删除对象
// @param minioClient
// @param bulkName 桶名
// @param objectName 对象名
// @return error
func (service *Service) RemoveObject(bulkName, objectName string) error {
	return service.ossService.(*minioService.MinioService).RemoveObject(bulkName, objectName)
}

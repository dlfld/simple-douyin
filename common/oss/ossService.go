package oss

// OSSService Oss 服务接口
type OSSService interface {
	// GetClient 获取连接客户端
	GetClient()
	// CreateBucket 创建桶
	CreateBucket(bucketName string)
	// DeleteBucket 删除桶
	DeleteBucket(bucketName string)
	// 上传文件
}

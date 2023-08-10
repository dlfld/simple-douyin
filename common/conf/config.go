package conf

var Mysql = mysqlConfig{
	Login: ".\/?charset=utf8&parseTime=True&loc=Local",
	Debug: true,
}
var Redis = redisConfig{
	Addr:     "42.192.46.30:6379",
	Password: "abc123456",
}
var UserService = userServiceConfig{
	Name: "userService",
	Addr: "127.0.0.1:8081",
}
var RelationService = relationServiceConfig{
	Name: "relationService",
	Addr: "127.0.0.1:8082",
}
var MessageService = messageServiceConfig{
	Name: "messageService",
	Addr: "127.0.0.1:8083",
}
var VideoService = videoServiceConfig{
	Name: "videoService",
	Addr: "127.0.0.1:8084",
}
var InteractionService = interactionServiceConfig{
	Name: "interactionService",
	Addr: "127.0.0.1:8085",
}

var MinioConfig = minioConfig{
	IP:                        "42.192.46.30",
	EndPoint:                  "42.192.46.30:9000",
	AccessKeyId:               "LX5CNH0ZL1I0BF6I4965",
	SecretAccessKey:           "75+9VGc4jBsQPzkJdvqgZeN6u6p3O+NnfF0KYxPY",
	UseSSL:                    false,
	VideoBucketName:           "video",
	AvatarBucketName:          "avatar",
	BackgroundImageBucketName: "bgi",
}

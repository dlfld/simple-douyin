package conf

var Mysql = mysqlConfig{
	Login: "root:abc123456@tcp(101.34.81.220:3306)/douyin?charset=utf8&parseTime=True&loc=Local",
	Debug: false,
}
var Redis = redisConfig{
	Addr:     "redis:6379",
	Password: "",
}

//TODO ip改成远端地址

var UserService = userServiceConfig{
	Name: "userService",
	Addr: "0.0.0.0:8081",
}
var RelationService = relationServiceConfig{
	Name: "relationService",
	Addr: "0.0.0.0:8082",
}
var MessageService = messageServiceConfig{
	Name: "messageService",
	Addr: "0.0.0.0:8083",
}
var VideoService = videoServiceConfig{
	Name: "videoService",
	Addr: "0.0.0.0:8084",
}
var InteractionService = interactionServiceConfig{
	Name: "interactionService",
	Addr: "0.0.0.0:8085",
}

var Kafka = kafkaConfig{
	//Addr: . "101.34.81.220",
	Addr: "10.23.65.200",
	Port: 9092,
}

var MinioConfig = minioConfig{
	IP:                        "101.34.81.220",
	EndPoint:                  "101.34.81.220:9000",
	AccessKeyId:               "LX5CNH0ZL1I0BF6I4965",
	SecretAccessKey:           "75+9VGc4jBsQPzkJdvqgZeN6u6p3O+NnfF0KYxPY",
	UseSSL:                    false,
	VideoBucketName:           "video",
	AvatarBucketName:          "avatar",
	BackgroundImageBucketName: "bgi",
}

func GetAllServiceName() []string {
	return []string{UserService.Name, RelationService.Name, MessageService.Name, VideoService.Name, InteractionService.Name}
}

var CosConfig = cosConfig{
	Url:       "https://douyin-1300206677.cos.ap-shanghai.myqcloud.com",
	SecretID:  "AKIDMdjfAWtCTvFcw794sgP2UvEOMrMgtz11",
	SecretKey: "6JRWbdZKzhW7TabR7eCcv1uc9hCcpxX6",
	ReginUrl:  "https://cos.COS_REGION.myqcloud.com",
}

var BloomConfig = bloomConfig{
	BloomBit: 1000000,
	HashNum:  3,
}

var EtcdConfig = etcdConfig{
	"http://etcd:2379",
}

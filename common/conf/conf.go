package conf

import "github.com/spf13/viper"

var (
	Mysql              mysqlConfig
	Redis              redisConfig
	UserService        userServiceConfig
	RelationService    relationServiceConfig
	MessageService     messageServiceConfig
	VideoService       videoServiceConfig
	InteractionService interactionServiceConfig
	Kafka              kafkaConfig
	CosConfig          cosConfig
	BloomConfig        bloomConfig
	EtcdConfig         etcdConfig
)

func init() {
	fileType := "yaml"
	//初始化mysql配置
	go func() {
		v := viper.New()
		v.SetConfigType(fileType) //设置配置文件类型
		v.SetConfigName("mysql")  //设置配置文件名
		v.AddConfigPath("config/mysql.yaml")
		Mysql = mysqlConfig{
			Login: v.Get("name").(string),
			Debug: v.GetBool("debug"),
		}
	}()
	//初始化redis配置
	go func() {
		v := viper.New()
		v.SetConfigType(fileType) //设置配置文件类型
		v.SetConfigName("redis")  //设置配置文件名
		v.AddConfigPath("config/redis.yaml")
		Redis = redisConfig{
			Addr:     v.Get("addr").(string),
			Password: v.Get("pwd").(string),
		}
	}()
	go func() {
		v := viper.New()
		v.SetConfigType(fileType) //设置配置文件类型
		v.SetConfigName("user")   //设置配置文件名
		v.AddConfigPath("config/user.yaml")
		UserService = userServiceConfig{
			Name: v.GetString("name"),
			Addr: v.GetString("addr"),
		}
	}()

	//初始化relationService配置
	go func() {
		v := viper.New()
		v.SetConfigType(fileType)   //设置配置文件类型
		v.SetConfigName("relation") //设置配置文件名
		v.AddConfigPath("config/relation.yaml")
		RelationService = relationServiceConfig{
			Name: v.GetString("name"),
			Addr: v.GetString("addr"),
		}
	}()

	//初始化messageService配置
	go func() {
		v := viper.New()
		v.SetConfigType(fileType)  //设置配置文件类型
		v.SetConfigName("message") //设置配置文件名
		v.AddConfigPath("config/message.yaml")
		MessageService = messageServiceConfig{
			Name: v.GetString("name"),
			Addr: v.GetString("addr"),
		}
	}()

	//初始化videoService配置
	go func() {
		v := viper.New()
		v.SetConfigType(fileType) //设置配置文件类型
		v.SetConfigName("video")  //设置配置文件名
		v.AddConfigPath("config/video.yaml")
		VideoService = videoServiceConfig{
			Name: v.GetString("name"),
			Addr: v.GetString("addr"),
		}
	}()

	//初始化interactionService配置
	go func() {
		v := viper.New()
		v.SetConfigType(fileType)      //设置配置文件类型
		v.SetConfigName("interaction") //设置配置文件名
		v.AddConfigPath("config/interaction.yaml")
		InteractionService = interactionServiceConfig{
			Name: v.GetString("name"),
			Addr: v.GetString("addr"),
		}
	}()

	//初始化kafka配置
	go func() {
		v := viper.New()
		v.SetConfigType(fileType) //设置配置文件类型
		v.SetConfigName("kafka")  //设置配置文件名
		v.AddConfigPath("config/kafka.yaml")
		Kafka = kafkaConfig{
			//Addr: . "101.34.81.220",
			Addr: v.GetString("addr"),
			Port: v.GetInt("port"),
		}
	}()

	//初始化cos配置
	go func() {
		v := viper.New()
		v.SetConfigType(fileType) //设置配置文件类型
		v.SetConfigName("cos")    //设置配置文件名
		v.AddConfigPath("config/cos.yaml")
		CosConfig = cosConfig{
			Url:       v.GetString("url"),
			SecretID:  v.GetString("secretId"),
			SecretKey: v.GetString("secretKey"),
			ReginUrl:  v.GetString("reginUrl"),
		}
	}()

	//初始化bloom配置
	go func() {
		v := viper.New()
		v.SetConfigType(fileType) //设置配置文件类型
		v.SetConfigName("bloom")  //设置配置文件名
		v.AddConfigPath("config/bloom.yaml")
		BloomConfig = bloomConfig{
			BloomBit: v.GetUint("bloomBit"),
			HashNum:  v.GetUint("hashNum"),
		}
	}()

	//初始化etcd配置
	go func() {
		v := viper.New()
		v.SetConfigType(fileType) //设置配置文件类型
		v.SetConfigName("etcd")   //设置配置文件名
		v.AddConfigPath("config/etcd.yaml")
		EtcdConfig = etcdConfig{
			v.GetString("addr"),
		}
	}()
}

// MinioConfig 废弃
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

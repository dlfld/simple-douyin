package conf

type mysqlConfig struct {
	Login string
	Debug bool
}

type redisConfig struct {
	Addr     string
	Password string
}

// minio配置文件
type minioConfig struct {
	IP                        string // ip
	EndPoint                  string // url
	AccessKeyId               string // key
	SecretAccessKey           string // password
	UseSSL                    bool   // is use ssl
	VideoBucketName           string // 存视频桶的名字
	AvatarBucketName          string // 存头像桶名字
	BackgroundImageBucketName string // 存背景图片桶名字
}

// cosConfig 访问各 API 所需的基础 URL
type cosConfig struct {
	Url       string //url
	SecretID  string
	SecretKey string
	ReginUrl  string
}
type userServiceConfig struct {
	Name string
	Addr string
}
type messageServiceConfig struct {
	Name string
	Addr string
}
type videoServiceConfig struct {
	Name string
	Addr string
}
type relationServiceConfig struct {
	Name string
	Addr string
}
type interactionServiceConfig struct {
	Name string
	Addr string
}

type kafkaConfig struct {
	Addr string
	Port int
}

type bloomConfig struct {
	Addr     string
	Password string
}

type etcdConfig struct {
	Addr string
}

type otelConfig struct {
	Addr string
}

type chatGptConfig struct {
	Name string
	Addr string
}

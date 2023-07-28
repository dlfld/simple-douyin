package conf

type mysqlConfig struct {
	Login string
	Debug bool
}

type redisConfig struct {
	Addr     string
	Password string
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

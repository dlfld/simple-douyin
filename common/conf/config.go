package conf

var Mysql = mysqlConfig{
	Login: "root:abc123456@tcp(42.192.46.30:3306)/douyin?charset=utf8&parseTime=True&loc=Local",
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

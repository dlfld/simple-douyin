package dao

type Dao struct {
	Mysql *mysql
	Redis *redis
}

func NewDao() (dao *Dao) {
	dao = &Dao{
		Mysql: NewMysql(),
		Redis: NewRedis(),
	}
	return dao
}

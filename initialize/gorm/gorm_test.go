package initialize

import "testing"

func TestMysqlConnect(t *testing.T) {
	GormInit()
	CreateTable()
}

package mysql

import (
	"fmt"
	"testing"
)

func TestNewMysqlConn(t *testing.T) {
	db, err := NewMysqlConn()
	if err != nil {
		t.Error(err.Error())
		return
	}
	tableNames := make([]string, 0)
	if err = db.Raw("show tables").Scan(&tableNames).Error; err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(tableNames)
}

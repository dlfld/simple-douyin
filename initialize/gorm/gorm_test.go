package initialize

import "testing"

func TestMysqlCreate(t *testing.T) {
	CreateTable()
}

func TestInteractionMysqlCreate(t *testing.T) {
	CreateInteractionTable()
}

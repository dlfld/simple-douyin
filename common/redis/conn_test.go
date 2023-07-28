package redis

import "testing"

func TestNewRedisConn(t *testing.T) {
	_, err := NewRedisConn()
	if err != nil {
		t.Error(err)
		return
	}
}

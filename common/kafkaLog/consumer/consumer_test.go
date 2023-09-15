package consumer

import (
	"fmt"
	"testing"
)

func TestReadLogFromKafka(t *testing.T) {
	for {
		key, value, err := PopLog()
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("key: %s  value: %v\n", key, value)
	}
}

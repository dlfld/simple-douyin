package main

import (
	"fmt"
	"testing"
)

func TestLogger(t *testing.T) {
	i := 0
	for i < 100 {
		logCollector.Info(fmt.Sprintf("%d", i))
	}

}

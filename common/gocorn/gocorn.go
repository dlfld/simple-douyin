package gocorn

import (
	"sync"
	"time"

	"github.com/go-co-op/gocron"
)

var once sync.Once
var s *gocron.Scheduler

func NewGocorn() *gocron.Scheduler {

	once.Do(
		func() {
			s = gocron.NewScheduler(time.Local)
			s.StartAsync()
		},
	)
	return s
}

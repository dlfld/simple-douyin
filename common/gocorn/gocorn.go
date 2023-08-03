package gocorn

import (
	"sync"
	"time"

	"github.com/go-co-op/gocron"
)

var once sync.Once

func NewGocorn() *gocron.Scheduler {
	var s *gocron.Scheduler
	once.Do(
		func() {
			s = gocron.NewScheduler(time.Local)
			s.StartAsync()
		},
	)
	return s
}

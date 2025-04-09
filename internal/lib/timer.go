package lib

import (
	"log"
	"time"
)

type Timer struct {
	start time.Time
}

func StartTimer() Timer {
	return Timer{
		start: time.Now(),
	}
}

func (t *Timer) End() {
	elapsed := time.Since(t.start)
	log.Printf("Operation took %s", elapsed)
}

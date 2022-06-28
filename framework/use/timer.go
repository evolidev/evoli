package use

import "time"

type Timer struct {
	start time.Time
}

func TimeRecord() *Timer {
	return &Timer{start: time.Now()}
}

func (t *Timer) Elapsed() time.Duration {
	return time.Now().Sub(t.start)
}

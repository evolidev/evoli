package use

import (
	"evoli.dev/framework/console/color"
	"time"
)

type timer struct {
	start time.Time
}

func TimeRecord() *timer {
	return &timer{start: time.Now()}
}

func (t *timer) Elapsed() time.Duration {
	return time.Now().Sub(t.start)
}

func (t *timer) ElapsedColored() string {
	return color.Text(150, "("+t.Elapsed().String()+")")
}

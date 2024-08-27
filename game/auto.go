package engine

import (
	"time"
)

// TimerObject represents an object with timer support.
type TimerObject struct {
	timer *time.Timer
}

// NewTimerObject creates a new TimerObject.
func NewTimerObject() *TimerObject {
	return &TimerObject{}
}

// StartTimer starts the timer with the specified duration and event handler.
func (t *TimerObject) StartTimer(duration time.Duration, eventHandler func()) {
	t.timer = time.AfterFunc(duration, eventHandler)
}

// StopTimer stops the timer.
func (t *TimerObject) StopTimer() {
	if t.timer != nil {
		t.timer.Stop()
	}
}

// FireEvent manually triggers the event handler.
func (t *TimerObject) FireEvent() {
	if t.timer != nil {
		t.timer.Reset(0)
	}
}

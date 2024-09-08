package engine

import (
	mylog "go-pk-server/log"
	"sync"
	"time"
)

// TimerObject represents an object with timer support.
type AutoInputProducer struct {
	timer         *time.Timer
	engineInputCh chan Input
	sync.Mutex
}

// NewTimerObject creates a new TimerObject.
func NewAutoInputProducer(engineInputCh chan Input) *AutoInputProducer {
	return &AutoInputProducer{
		engineInputCh: engineInputCh,
	}
}

// StartTimer starts the timer with the specified duration and event handler.
func (aip *AutoInputProducer) CreatGameInputAfter(input GaneInputType, duration time.Duration) {
	// lock the timer
	aip.Lock()
	defer aip.Unlock()

	aip.timer = time.AfterFunc(duration, func() {
		input := Input{Type: input}
		aip.engineInputCh <- input
		mylog.Infof("Auto input [ %v] to the GAME engine", input)
	})
	// Log the timer creation
	mylog.Infof("Will auto produce an input %v after %v", input, duration)
}

// StopTimer stops the timer.
func (aip *AutoInputProducer) StopOngoingAutoInput() {
	// lock the timer
	aip.Lock()
	defer aip.Unlock()

	if aip.timer != nil {
		aip.timer.Stop()
	}
}

// FireEvent manually triggers the event handler.
func (aip *AutoInputProducer) ResetAutoInputDuration() {
	// lock the timer
	aip.Lock()
	defer aip.Unlock()

	if aip.timer != nil {
		aip.timer.Reset(0)
	}
}

package esc

import (
	"time"

	"github.com/marksaravi/drone-go/models"
)

const (
	NUM_OF_ESC uint8 = 4
)

type powerbreaker interface {
	Connect()
	Disconnect()
}

type pwmDevice interface {
	SetThrottles(map[int]float64)
}

type escDev struct {
	pwmDev         pwmDevice
	powerbreaker   powerbreaker
	lastUpdate     time.Time
	updateInterval time.Duration
	debug          bool
}

func NewESC(pwmDev pwmDevice, powerbreaker powerbreaker, updatesPerSecond int, debug bool) *escDev {
	powerbreaker.Disconnect()
	return &escDev{
		pwmDev:         pwmDev,
		powerbreaker:   powerbreaker,
		lastUpdate:     time.Now().Add(-time.Second),
		updateInterval: time.Second / time.Duration(updatesPerSecond),
		debug:          debug,
	}
}

func (e *escDev) zeroThrottle() {
	e.pwmDev.SetThrottles(map[int]float64{0: 0, 1: 0, 2: 0, 3: 0})
}

func (e *escDev) On() {
	e.zeroThrottle()
	if !e.debug {
		e.powerbreaker.Connect()
	}
}

func (e *escDev) Off() {
	e.zeroThrottle()
	e.powerbreaker.Disconnect()
}

func (e *escDev) SetThrottles(throttles models.Throttles) {
	if time.Since(e.lastUpdate) >= e.updateInterval {
		e.lastUpdate = time.Now()
		func(throttles map[int]float64) {
			e.pwmDev.SetThrottles(throttles)
		}(throttles.Throttles)
	}
}

package esc

import (
	"time"
)

const (
	NUM_OF_ESC uint8 = 4
)

type powerbreaker interface {
	Connect()
	Disconnect()
}

type pwmDevice interface {
	SetThrottles([]float64) error
	NumberOfChannels()      int

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
	z:=make([]float64, e.pwmDev.NumberOfChannels())
	for i:=0; i<len(z); i++ {
		z[i]=0
	}
	e.pwmDev.SetThrottles(z)
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

func (e *escDev) SetThrottles(throttles []float64) {
	if time.Since(e.lastUpdate) >= e.updateInterval {
		e.lastUpdate = time.Now()
		func(throttles []float64) {
			e.pwmDev.SetThrottles(throttles)
		}(throttles)
	}
}

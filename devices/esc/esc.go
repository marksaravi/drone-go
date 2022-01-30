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
	OffAll()
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

func (e *escDev) On() {
	e.pwmDev.OffAll()
	if !e.debug {
		e.powerbreaker.Connect()
	}
}

func (e *escDev) Off() {
	e.pwmDev.OffAll()
	e.powerbreaker.Disconnect()
}

func (e *escDev) SetThrottles(throttles models.Throttles) {
	if throttles.Active {
		if time.Since(e.lastUpdate) >= e.updateInterval {
			e.lastUpdate = time.Now()
			e.pwmDev.SetThrottles(throttles.Throttles)
		}

	} else {
		e.pwmDev.OffAll()
	}
}

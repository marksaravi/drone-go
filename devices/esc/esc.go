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
	NumberOfChannels() int
}

type escDev struct {
	pwmDev                   pwmDevice
	motorsToChannelsMappings []int
	throttles                []float64
	powerbreaker             powerbreaker
	lastUpdate               time.Time
	updateInterval           time.Duration
	debug                    bool
}

func NewESC(
	pwmDev pwmDevice,
	motorsToChannelsMappings []int,
	powerbreaker powerbreaker,
	updatesPerSecond int,
	debug bool,
) *escDev {
	powerbreaker.Disconnect()
	mappings := make([]int, 4)
	copy(mappings, motorsToChannelsMappings)
	return &escDev{
		pwmDev:                   pwmDev,
		motorsToChannelsMappings: mappings,
		throttles:                make([]float64, pwmDev.NumberOfChannels()),
		powerbreaker:             powerbreaker,
		lastUpdate:               time.Now().Add(-time.Second),
		updateInterval:           time.Second / time.Duration(updatesPerSecond),
		debug:                    debug,
	}
}

func (e *escDev) zeroThrottle() {
	z := make([]float64, e.pwmDev.NumberOfChannels())
	for i := 0; i < len(z); i++ {
		z[i] = 0
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

func (e *escDev) mapMotorsToChannels(motors []float64) {
	for i := 0; i < len(e.throttles); i++ {
		e.throttles[i] = 0
	}
	for i := 0; i < len(motors); i++ {
		e.throttles[e.motorsToChannelsMappings[i]] = motors[i]
	}
}

func (e *escDev) SetThrottles(motors []float64) {
	e.mapMotorsToChannels(motors)
	if time.Since(e.lastUpdate) >= e.updateInterval {
		e.lastUpdate = time.Now()
		func() {
			e.pwmDev.SetThrottles(e.throttles)
		}()
	}
}

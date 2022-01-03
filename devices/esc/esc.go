package esc

import (
	"log"
	"sync"

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
	SetThrottle(int, float32)
}

type escDev struct {
	pwmDev                 pwmDevice
	powerbreaker           powerbreaker
	pwmDeviceToESCMappings map[int]int
	throttels              models.Throttles
	throttlesChan          chan models.Throttles
	isActive               bool
}

func NewESC(pwmDev pwmDevice, powerbreaker powerbreaker, pwmDeviceToESCMappings map[int]int) *escDev {
	return &escDev{
		pwmDev:                 pwmDev,
		powerbreaker:           powerbreaker,
		pwmDeviceToESCMappings: pwmDeviceToESCMappings,
		throttlesChan:          make(chan models.Throttles),
		isActive:               true,
	}
}

func (e *escDev) On() {
	e.offAll()
	e.powerbreaker.Connect()
}

func (e *escDev) Off() {
	e.powerbreaker.Disconnect()
	e.offAll()
}

func (e *escDev) SetThrottles(throttles models.Throttles) {
	if e.isActive {
		e.throttlesChan <- throttles
	}
}

func (e *escDev) SetThrottle(channel uint8, throttle float64) {
	if e.isActive {
		throttels := e.throttels
		throttels[channel] = throttle
		e.throttlesChan <- throttels
	}
}

func (e *escDev) Close() {
	if e.isActive {
		close(e.throttlesChan)
	}
}

func (e *escDev) Start(wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer log.Println("ESC is closed")

		for e.isActive {
			throttels, ok := <-e.throttlesChan
			if ok {
				var ch uint8
				for ch = 0; ch < NUM_OF_ESC; ch++ {
					e.pwmDev.SetThrottle(e.pwmDeviceToESCMappings[int(ch)], float32(throttels[ch]))
				}
			} else {
				e.isActive = false
			}
		}
	}()
}

func (e *escDev) offAll() {
	if e.isActive {
		e.throttlesChan <- models.Throttles{0: 0, 1: 0, 2: 0, 3: 0}
	}
}

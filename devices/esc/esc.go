package esc

import (
	"log"
	"sync"
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
	SetThrottle(int, float32)
}

type escDev struct {
	pwmDev                 pwmDevice
	powerbreaker           powerbreaker
	pwmDeviceToESCMappings map[int]int
	throttels              models.Throttles
	throttlesChan          chan models.Throttles
	lastUpdate             time.Time
	updateInterval         time.Duration
	isActive               bool
	debug                  bool
}

func NewESC(pwmDev pwmDevice, powerbreaker powerbreaker, updatesPerSecond int, pwmDeviceToESCMappings map[int]int, debug bool) *escDev {
	return &escDev{
		pwmDev:                 pwmDev,
		powerbreaker:           powerbreaker,
		pwmDeviceToESCMappings: pwmDeviceToESCMappings,
		throttlesChan:          make(chan models.Throttles),
		lastUpdate:             time.Now(),
		updateInterval:         time.Second / time.Duration(updatesPerSecond),
		isActive:               true,
		debug:                  debug,
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
	if e.isActive && time.Since(e.lastUpdate) >= e.updateInterval {
		e.lastUpdate = time.Now()
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
	log.Println("Starting ESC...")

	go func() {
		defer wg.Done()
		defer log.Println("ESC is closed")

		for e.isActive {
			throttels, ok := <-e.throttlesChan
			if ok {
				showThrottles(throttels)
				if !e.debug {
					var ch uint8
					for ch = 0; ch < NUM_OF_ESC; ch++ {
						e.pwmDev.SetThrottle(e.pwmDeviceToESCMappings[int(ch)], float32(throttels[ch]))
					}
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

var lastPrint time.Time = time.Now()

func showThrottles(throttles models.Throttles) {
	if time.Since(lastPrint) > time.Second/2 {
		lastPrint = time.Now()
		log.Printf("0: %6.2f, 1: %6.2f, 2: %6.2f, 3: %6.2f\n", throttles[0], throttles[1], throttles[2], throttles[3])
	}
}

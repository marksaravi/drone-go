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
	SetThrottles(models.Throttles)
}

type escDev struct {
	pwmDev         pwmDevice
	powerbreaker   powerbreaker
	throttlesChan  chan models.Throttles
	lastUpdate     time.Time
	updateInterval time.Duration
	isActive       bool
	debug          bool
}

func NewESC(pwmDev pwmDevice, powerbreaker powerbreaker, updatesPerSecond int, debug bool) *escDev {
	return &escDev{
		pwmDev:         pwmDev,
		powerbreaker:   powerbreaker,
		throttlesChan:  make(chan models.Throttles),
		lastUpdate:     time.Now(),
		updateInterval: time.Second / time.Duration(updatesPerSecond),
		isActive:       true,
		debug:          debug,
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
				// showThrottles(throttels)
				if !e.debug {
					e.pwmDev.SetThrottles(throttels)
				}
			} else {
				e.isActive = false
			}
		}
	}()
}

func (e *escDev) offAll() {
	if e.isActive {
		e.throttlesChan <- models.Throttles{
			Throttle:         0,
			ControlVariables: map[int]float64{0: 0, 1: 0, 2: 0, 3: 0},
		}
	}
}

// var lastPrint time.Time = time.Now()

// func showThrottles(throttles models.Throttles) {
// 	if time.Since(lastPrint) > time.Second/2 {
// 		lastPrint = time.Now()
// 		log.Printf("0: %6.2f, 1: %6.2f, 2: %6.2f, 3: %6.2f\n", throttles[0], throttles[1], throttles[2], throttles[3])
// 	}
// }

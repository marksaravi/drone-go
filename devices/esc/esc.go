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
	if !e.isActive {
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
					e.pwmDev.SetThrottle(int(ch), float32(throttels[ch]))
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

// func (escdev *escDev) setThrottle(motor uint8, throttle float32) {
// 	e.esc.SetThrottle(e.motorsEscMappings[int(motor)], throttle)
// }

// func (escdev *escDev) allMotorsOff() {
// 	var i uint8
// 	for i = 0; i < NUM_OF_ESC; i++ {
// 		mc.setThrottle(i, 0)
// 	}
// }

// func NewESC() *escDev {
// 	flightControlConfigs := config.ReadConfigs().FlightControl
// 	escConfigs := flightControlConfigs.ESC
// 	powerbreaker := hardware.NewGPIOOutput(flightControlConfigs.PowerBreaker)
// 	b, _ := i2creg.Open(escConfigs.I2CDev)
// 	i2cConn := &i2c.Dev{Addr: pca9685.PCA9685Address, Bus: b}
// 	pwmDev, err := pca9685.NewPCA9685(pca9685.PCA9685Address, i2cConn, escConfigs.MaxThrottle)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	esc := NewMotorsControl(pwmDev, powerbreaker, escConfigs.PwmDeviceToESCMappings)
// 	return esc
// }

// func NewThrottleChannel(wg *sync.WaitGroup) (chan<- models.Throttles, chan<- bool) {
// 	wg.Add(1)
// 	throttleChannel := make(chan map[uint8]float64)
// 	onOff := make(chan bool)
// 	esc := NewESC()
// 	go escRoutine(wg, esc, throttleChannel, onOff)
// 	return throttleChannel, onOff
// }

// func escRoutine(wg *sync.WaitGroup, esc *escDev, throttleChannel <-chan models.Throttles, onOff <-chan bool) {
// 	defer wg.Done()
// 	defer log.Println("ESC stopped")
// 	for throttleChannel != nil || onOff != nil {
// 		select {
// 		case throttles, ok := <-throttleChannel:
// 			if ok {
// 				esc.SetThrottles(throttles)
// 			} else {
// 				throttleChannel = nil
// 			}
// 		case on, ok := <-onOff:
// 			if ok {
// 				if on {
// 					esc.On()
// 				} else {
// 					esc.Off()
// 				}
// 			} else {
// 				esc.Off()
// 				onOff = nil
// 			}
// 		}
// 	}
// }

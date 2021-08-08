package esc

import (
	"fmt"
	"os"
	"time"

	"github.com/MarkSaravi/drone-go/connectors/i2c"
	"github.com/MarkSaravi/drone-go/hardware/pca9685"
	"github.com/MarkSaravi/drone-go/modules/powerbreaker"
	"github.com/MarkSaravi/drone-go/types"
)

type esc interface {
	Start() error
	SetThrottle(motor int, throttle float32)
	StopAll()
	Close()
}

type breaker interface {
	MotorsOn()
	MotorsOff()
}

type escsHandler struct {
	esc            esc
	breaker        breaker
	throttles      []types.Throttle
	lastUpdate     time.Time
	updateInterval time.Duration
}

func (h *escsHandler) MotorsOn() {
	for i := 0; i < 4; i++ {
		h.esc.SetThrottle(i, 0)
	}
	h.breaker.MotorsOn()
	time.Sleep(time.Millisecond * 2000)
}

func (h *escsHandler) MotorsOff() {
	h.esc.StopAll()
	h.breaker.MotorsOff()
	h.esc.Close()
}

func (h *escsHandler) SetThrottles(throttles []types.Throttle) {
	h.throttles = throttles
	if time.Since(h.lastUpdate) < h.updateInterval {
		return
	}
	h.lastUpdate = time.Now()
	go func() {
		for _, throttle := range h.throttles {
			h.esc.SetThrottle(throttle.Motor, throttle.Value)
		}
	}()
}

func NewESCsHandler(config types.EscConfig) *escsHandler {
	i2cConnection, err := i2c.Open(config.SPIDevice)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	breaker := powerbreaker.NewPowerBreaker(config.PowerBrokerGPIO)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	pwmDev, err := pca9685.NewPCA9685Driver(pca9685.PCA9685Address, i2cConnection, config.MaxThrottle, map[int]int{0: 0, 1: 1, 2: 2, 3: 4})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = pwmDev.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return &escsHandler{
		esc:            pwmDev,
		breaker:        breaker,
		lastUpdate:     time.Now(),
		updateInterval: time.Second / time.Duration(config.UpdateFrequency),
	}
}

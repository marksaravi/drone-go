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
	esc       esc
	breaker   breaker
	throttles []types.Throttle
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
	for _, throttle := range h.throttles {
		h.esc.SetThrottle(throttle.Motor, throttle.Value)
	}
}

func NewESCsHandler() *escsHandler {
	i2cConnection, err := i2c.Open("/dev/i2c-1")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	breaker := powerbreaker.NewPowerBreaker()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	pwmDev, err := pca9685.NewPCA9685Driver(pca9685.PCA9685Address, i2cConnection)
	err = pwmDev.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return &escsHandler{
		esc:     pwmDev,
		breaker: breaker,
	}
}

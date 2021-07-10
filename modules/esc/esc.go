package esc

import (
	"fmt"
	"os"

	"github.com/MarkSaravi/drone-go/connectors/i2c"
	"github.com/MarkSaravi/drone-go/hardware/pca9685"
	"github.com/MarkSaravi/drone-go/modules/powerbreaker"
	"github.com/MarkSaravi/drone-go/types"
)

type esc interface {
	Start(frequency float32) error
	SetPulseWidth(channel int, pulseWidth float32)
	SetPulseWidthAll(pulseWidth float32)
	StopAll()
	Halt() error
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
	h.breaker.MotorsOn()
}

func (h *escsHandler) MotorsOff() {
	h.breaker.MotorsOff()
}

func (h *escsHandler) SetThrottles(throttles []types.Throttle) {
	h.throttles = throttles
	// for _, throttle := range h.throttles {
	// 	h.esc.SetPulseWidth(throttle.Motor, throttle.Value)
	// }
}

func NewESCsHandler() *escsHandler {
	i2cConnection, err := i2c.Open("/dev/i2c-1")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	breaker := powerbreaker.NewPowerBreaker()
	pwmDev, err := pca9685.NewPCA9685Driver(pca9685.PCA9685Address, i2cConnection)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return &escsHandler{
		esc:     pwmDev,
		breaker: breaker,
	}
}

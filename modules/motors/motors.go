package motors

import (
	"github.com/MarkSaravi/drone-go/types"
)

const (
	NUM_OF_MOTORS int = 4
)

type motorsControl struct {
	esc          types.ESC
	powerbreaker types.PowerBreaker
}

func NewMotorsControl(esc types.ESC, powerbreaker types.PowerBreaker) types.MotorsController {
	return &motorsControl{
		esc:          esc,
		powerbreaker: powerbreaker,
	}
}

func (mc *motorsControl) On() {
	mc.allMotorsOff()
	mc.powerbreaker.Connect()
}

func (mc *motorsControl) Off() {
	mc.allMotorsOff()
	mc.powerbreaker.Disconnect()
}

func (mc *motorsControl) SetThrottles(throttles map[int]float32) {
	for motor := 0; motor < NUM_OF_MOTORS; motor++ {
		mc.esc.SetThrottle(motor, throttles[motor])
	}
}

func (mc *motorsControl) allMotorsOff() {
	for motor := 0; motor < NUM_OF_MOTORS; motor++ {
		mc.esc.SetThrottle(motor, 0)
	}
}

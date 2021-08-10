package motors

import (
	"github.com/MarkSaravi/drone-go/modules/powerbreaker"
)

const (
	NUM_OF_MOTORS int = 4
)

type Throttle struct {
	Motor int
	Value float32
}

type ESC interface {
	SetThrottle(int, float32)
}

type MotorsController interface {
	On()
	Off()
	SetThrottles(map[int]float32)
}

type motorsControl struct {
	esc          ESC
	powerbreaker powerbreaker.PowerBreaker
}

func NewMotorsControl(esc ESC, powerbreaker powerbreaker.PowerBreaker) MotorsController {
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

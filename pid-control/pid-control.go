package pidcontrol

import (
	"github.com/MarkSaravi/drone-go/models"
)

type pidControl struct {
	commands  models.FlightCommands
	rotations models.ImuRotations
	throttles map[uint8]float32
}

func NewPIDControl() *pidControl {
	return &pidControl{}
}

func (pid *pidControl) ApplyFlightCommands(flightCommands models.FlightCommands) {
	pid.commands = flightCommands
	if t, err := pid.calcMotorsThrottles(); err == nil {
		pid.throttles = t
	}
}

func (pid *pidControl) ApplyRotations(rotations models.ImuRotations) {
	pid.rotations = rotations
	if t, err := pid.calcMotorsThrottles(); err == nil {
		pid.throttles = t
	}
}

func (pid *pidControl) calcMotorsThrottles() (map[uint8]float32, error) {
	throttle := pid.commands.Throttle / 5 * 10
	return map[uint8]float32{
		0: throttle,
		1: throttle,
		2: throttle,
		3: throttle,
	}, nil
}

func (pid *pidControl) Throttles() map[uint8]float32 {
	return pid.throttles
}

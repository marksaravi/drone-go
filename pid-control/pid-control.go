package pidcontrol

import (
	"github.com/marksaravi/drone-go/config"
	"github.com/marksaravi/drone-go/models"
)

type pidConfig struct {
	pGain                 float64
	iGain                 float64
	dGain                 float64
	analogInputToThrottle float64
}

type pidControl struct {
	analogInputToThrottle float64
	pGain                 float64
	iGain                 float64
	dGain                 float64
	commands              models.FlightCommands
	rotations             models.ImuRotations
	prevRotations         models.ImuRotations
	throttle              float64
	iThrottle             float64
	throttles             map[uint8]float32
}

func NewPIDControl() *pidControl {
	configs := config.ReadFlightControlConfig().Configs.PID
	return newPID(pidConfig{
		analogInputToThrottle: configs.AnalogInputToThrottle,
		pGain:                 configs.PGain,
		iGain:                 configs.IGain,
		dGain:                 configs.DGain,
	})
}

func (pid *pidControl) ApplyFlightCommands(flightCommands models.FlightCommands) {
	pid.commands = flightCommands
	if t, err := pid.calcMotorsThrottles(); err == nil {
		pid.throttles = t
	}
}

func (pid *pidControl) ApplyRotations(rotations models.ImuRotations) {
	pid.prevRotations = pid.rotations
	pid.rotations = rotations
	if t, err := pid.calcMotorsThrottles(); err == nil {
		pid.throttles = t
	}
}

func (pid *pidControl) calcMotorsThrottles() (map[uint8]float32, error) {
	throttle := pid.commands.Throttle * float32(pid.analogInputToThrottle)
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

func newPID(config pidConfig) *pidControl {
	return &pidControl{
		analogInputToThrottle: config.analogInputToThrottle,
		pGain:                 config.pGain,
		iGain:                 config.iGain,
		dGain:                 config.dGain,
	}
}

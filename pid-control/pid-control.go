package pidcontrol

import (
	"github.com/marksaravi/drone-go/config"
	"github.com/marksaravi/drone-go/models"
)

type pidConfig struct {
	pGain        float64
	iGain        float64
	dGain        float64
	throttleGain float64
}

type pidControl struct {
	commands     models.FlightCommands
	rotations    models.ImuRotations
	throttles    map[uint8]float32
	pGain        float64
	iGain        float64
	dGain        float64
	throttleGain float64
}

func NewPIDControl() *pidControl {
	configs := config.ReadFlightControlConfig().Configs.PID
	return newPID(pidConfig{
		pGain:        configs.PGain,
		iGain:        configs.IGain,
		dGain:        configs.DGain,
		throttleGain: configs.ThrottleGain,
	})
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
	throttle := pid.commands.Throttle * float32(pid.throttleGain)
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
		throttleGain: config.throttleGain,
		pGain:        config.pGain,
		iGain:        config.iGain,
		dGain:        config.dGain,
	}
}

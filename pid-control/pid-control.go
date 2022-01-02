package pidcontrol

import (
	"github.com/marksaravi/drone-go/config"
	"github.com/marksaravi/drone-go/constants"
	"github.com/marksaravi/drone-go/models"
)

const ROTATIONS_HISTORY_SIZE int = 2

type pidState struct {
	roll     float64
	pitch    float64
	yaw      float64
	throttle float64
}

type pidControl struct {
	pGain            float64
	iGain            float64
	dGain            float64
	state            pidState
	targetState      pidState
	rotationsHistory []models.ImuRotations
	throttles        models.Throttles
}

func NewPIDControl() *pidControl {
	configs := config.ReadConfigs().FlightControl.PID
	return &pidControl{
		pGain:            configs.PGain,
		iGain:            configs.IGain,
		dGain:            configs.DGain,
		rotationsHistory: make([]models.ImuRotations, ROTATIONS_HISTORY_SIZE),
		throttles:        models.Throttles{0: 0, 1: 0, 2: 0, 3: 0},
	}
}

func (pid *pidControl) SetFlightCommands(flightCommands models.FlightCommands) models.Throttles {
	pid.targetState = flightControlCommandToPIDCommand(flightCommands)
	pid.calcThrottles()
	return pid.throttles
}

func (pid *pidControl) SetRotations(rotations models.ImuRotations) models.Throttles {
	for i := 1; i < ROTATIONS_HISTORY_SIZE; i++ {
		pid.rotationsHistory[i] = pid.rotationsHistory[i-1]
	}
	pid.rotationsHistory[0] = rotations
	pid.calcThrottles()
	return pid.throttles
}

func (pid *pidControl) calcThrottles() {
	t := pid.state.throttle
	pid.throttles = models.Throttles{
		0: t,
		1: t,
		2: t,
		3: t,
	}
}

func (pid *pidControl) Throttles() models.Throttles {
	return pid.throttles
}

func flightControlCommandToPIDCommand(c models.FlightCommands) pidState {
	midValue := float64(constants.JOYSTICK_RESOLUTION / 2)

	return pidState{
		roll:     float64(c.Roll) - midValue,
		pitch:    float64(c.Pitch) - midValue,
		yaw:      float64(c.Yaw) - midValue,
		throttle: float64(c.Throttle),
	}
}

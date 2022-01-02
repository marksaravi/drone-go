package pidcontrol

import (
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
	maxRoll          float64
	maxPitch         float64
	maxYaw           float64
	maxThrottle      float64
	throttles        models.Throttles
	emergencyStop    bool
}

func NewPIDControl(pGain, iGain, dGain, maxRoll, maxPitch, maxYaw, maxThrottle float64) *pidControl {
	return &pidControl{
		pGain:            pGain,
		iGain:            iGain,
		dGain:            dGain,
		maxRoll:          maxRoll,
		maxPitch:         maxPitch,
		maxYaw:           maxYaw,
		maxThrottle:      maxThrottle,
		rotationsHistory: make([]models.ImuRotations, ROTATIONS_HISTORY_SIZE),
		throttles:        models.Throttles{0: 0, 1: 0, 2: 0, 3: 0},
		emergencyStop:    false,
	}
}

func (pid *pidControl) SetFlightCommands(flightCommands models.FlightCommands) {
	pid.targetState = flightControlCommandToPIDCommand(flightCommands)
	pid.calcThrottles()
}

func (pid *pidControl) SetRotations(rotations models.ImuRotations) {
	for i := 1; i < ROTATIONS_HISTORY_SIZE; i++ {
		pid.rotationsHistory[i] = pid.rotationsHistory[i-1]
	}
	pid.rotationsHistory[0] = rotations
	pid.calcThrottles()
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
	if pid.emergencyStop {
		pid.applyEmergencyStop()
	}

	return pid.throttles
}

func (pid *pidControl) SetEmergencyStop(stop bool) {
	pid.emergencyStop = stop
}

func (pid *pidControl) applyEmergencyStop() {
	pid.targetState = pidState{
		roll:     0,
		pitch:    0,
		yaw:      0,
		throttle: 0,
	}
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

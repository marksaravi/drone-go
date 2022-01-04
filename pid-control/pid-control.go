package pidcontrol

import (
	"github.com/marksaravi/drone-go/models"
)

type pidTargetState struct {
	roll     float64
	pitch    float64
	yaw      float64
	throttle float64
}

type pidState struct {
	roll      float64
	pitch     float64
	yaw       float64
	throttles models.Throttles
}

type pidControl struct {
	imuDataPerSecond        int
	imuDataBufferSize       int
	pGain                   float64
	iGain                   float64
	dGain                   float64
	targetState             pidTargetState
	state                   pidState
	rotationsHistory        []models.ImuRotations
	maxJoystickDigitalValue float64
	maxRoll                 float64
	maxPitch                float64
	maxYaw                  float64
	maxThrottle             float64
	emergencyStop           bool
}

func NewPIDControl(imuDataPerSecond int, pGain, iGain, dGain, maxRoll, maxPitch, maxYaw, maxThrottle float64, maxJoystickDigitalValue uint16) *pidControl {
	imuDataBufferSize := 2
	return &pidControl{
		imuDataPerSecond:        imuDataPerSecond,
		imuDataBufferSize:       imuDataBufferSize,
		pGain:                   pGain,
		iGain:                   iGain,
		dGain:                   dGain,
		maxRoll:                 maxRoll,
		maxPitch:                maxPitch,
		maxYaw:                  maxYaw,
		maxThrottle:             maxThrottle,
		maxJoystickDigitalValue: float64(maxJoystickDigitalValue),
		rotationsHistory:        make([]models.ImuRotations, imuDataBufferSize),
		targetState: pidTargetState{
			roll:     0,
			pitch:    0,
			yaw:      0,
			throttle: 0,
		},
		state: pidState{
			roll:      0,
			pitch:     0,
			yaw:       0,
			throttles: models.Throttles{0: 0, 1: 0, 2: 0, 3: 0},
		},
		emergencyStop: false,
	}
}

func (pid *pidControl) SetFlightCommands(flightCommands models.FlightCommands) {
	pid.targetState = pid.flightControlCommandToPIDCommand(flightCommands)
	pid.calcThrottles()
}

func (pid *pidControl) SetRotations(rotations models.ImuRotations) {
	for i := 1; i < pid.imuDataBufferSize; i++ {
		pid.rotationsHistory[i] = pid.rotationsHistory[i-1]
	}
	pid.rotationsHistory[0] = rotations
	pid.calcThrottles()
}

func (pid *pidControl) calcThrottles() {
	t := pid.targetState.throttle
	pid.state = pidState{
		throttles: models.Throttles{
			0: t,
			1: t,
			2: t,
			3: t,
		},
	}
}

func (pid *pidControl) Throttles() models.Throttles {
	return pid.state.throttles
}

func (pid *pidControl) SetEmergencyStop(stop bool) {
	pid.emergencyStop = stop
}

// func (pid *pidControl) applyEmergencyStop() {
// 	pid.targetState = pidTargetState{
// 		roll:     0,
// 		pitch:    0,
// 		yaw:      0,
// 		throttle: 0,
// 	}
// }

func (pid *pidControl) joystickToPidValue(joystickDigitalValue uint16, maxValue float64) float64 {
	normalizedDigitalValue := float64(joystickDigitalValue) - pid.maxJoystickDigitalValue/2
	return normalizedDigitalValue / pid.maxJoystickDigitalValue * maxValue
}

func (pid *pidControl) throttleToPidThrottle(joystickDigitalValue uint16) float64 {
	return float64(joystickDigitalValue) / pid.maxJoystickDigitalValue * pid.maxThrottle
}

func (pid *pidControl) flightControlCommandToPIDCommand(c models.FlightCommands) pidTargetState {

	return pidTargetState{
		roll:     pid.joystickToPidValue(c.Roll, pid.maxRoll),
		pitch:    pid.joystickToPidValue(c.Pitch, pid.maxPitch),
		yaw:      pid.joystickToPidValue(c.Yaw, pid.maxYaw),
		throttle: pid.throttleToPidThrottle(c.Throttle),
	}
}

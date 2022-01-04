package pid

import (
	"log"
	"time"

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

type pidControls struct {
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

func NewPIDControls(imuDataPerSecond int, pGain, iGain, dGain, maxRoll, maxPitch, maxYaw, maxThrottle float64, maxJoystickDigitalValue uint16) *pidControls {
	imuDataBufferSize := 2
	return &pidControls{
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

func (pidcontrols *pidControls) SetFlightCommands(flightCommands models.FlightCommands) {
	pidcontrols.targetState = pidcontrols.flightControlCommandToPIDCommand(flightCommands)
	showStates(pidcontrols.targetState)
	pidcontrols.calcThrottles()
}

func (pidcontrols *pidControls) SetRotations(rotations models.ImuRotations) {
	for i := 1; i < pidcontrols.imuDataBufferSize; i++ {
		pidcontrols.rotationsHistory[i] = pidcontrols.rotationsHistory[i-1]
	}
	pidcontrols.rotationsHistory[0] = rotations
	pidcontrols.calcThrottles()
}

func (pidcontrols *pidControls) calcThrottles() {
	t := pidcontrols.targetState.throttle
	pidcontrols.state = pidState{
		throttles: models.Throttles{
			0: t,
			1: t,
			2: t,
			3: t,
		},
	}
}

func (pidcontrols *pidControls) Throttles() models.Throttles {
	return pidcontrols.state.throttles
}

func (pidcontrols *pidControls) SetEmergencyStop(stop bool) {
	pidcontrols.emergencyStop = stop
}

// func (pidcontrols *pidControls) applyEmergencyStop() {
// 	pidcontrols.targetState = pidTargetState{
// 		roll:     0,
// 		pitch:    0,
// 		yaw:      0,
// 		throttle: 0,
// 	}
// }

func (pidcontrols *pidControls) joystickToPidValue(joystickDigitalValue uint16, maxValue float64) float64 {
	normalizedDigitalValue := float64(joystickDigitalValue) - pidcontrols.maxJoystickDigitalValue/2
	return normalizedDigitalValue / pidcontrols.maxJoystickDigitalValue * maxValue
}

func (pidcontrols *pidControls) throttleToPidThrottle(joystickDigitalValue uint16) float64 {
	return float64(joystickDigitalValue) / pidcontrols.maxJoystickDigitalValue * pidcontrols.maxThrottle
}

func (pidcontrols *pidControls) flightControlCommandToPIDCommand(c models.FlightCommands) pidTargetState {

	return pidTargetState{
		roll:     pidcontrols.joystickToPidValue(c.Roll, pidcontrols.maxRoll),
		pitch:    pidcontrols.joystickToPidValue(c.Pitch, pidcontrols.maxPitch),
		yaw:      pidcontrols.joystickToPidValue(c.Yaw, pidcontrols.maxYaw),
		throttle: pidcontrols.throttleToPidThrottle(c.Throttle),
	}
}

var lastPrint time.Time = time.Now()

func showStates(s pidTargetState) {
	if time.Since(lastPrint) > time.Second/2 {
		lastPrint = time.Now()
		log.Printf("roll: %6.2f, pitch: %6.2f, yaw: %6.2f, throttle: %6.2f\n", s.roll, s.pitch, s.yaw, s.throttle)
	}
}

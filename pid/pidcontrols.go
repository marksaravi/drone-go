package pid

import (
	"log"
	"time"

	"github.com/marksaravi/drone-go/models"
)

type pidState struct {
	roll     float64
	pitch    float64
	yaw      float64
	throttle float64
}

type pidControls struct {
	roll                    *pidControl
	pitch                   *pidControl
	yaw                     *pidControl
	targetState             pidState
	state                   pidState
	throttles               models.Throttles
	maxJoystickDigitalValue float64
	throttleLimit           float64
	axisAlignmentAngle      float64
	emergencyStop           bool
}

func NewPIDControls(
	pGain, iGain, dGain float64,
	limitRoll, limitPitch, limitYaw, throttleLimit float64,
	maxJoystickDigitalValue uint16,
	axisAlignmentAngle float64,
) *pidControls {

	return &pidControls{
		roll:                    NewPIDControl(pGain, iGain, dGain, limitRoll),
		pitch:                   NewPIDControl(pGain, iGain, dGain, limitPitch),
		yaw:                     NewPIDControl(pGain, iGain, dGain, limitYaw),
		throttleLimit:           throttleLimit,
		maxJoystickDigitalValue: float64(maxJoystickDigitalValue),
		axisAlignmentAngle:      axisAlignmentAngle,
		targetState: pidState{
			roll:     0,
			pitch:    0,
			yaw:      0,
			throttle: 0,
		},
		state: pidState{
			roll:     0,
			pitch:    0,
			yaw:      0,
			throttle: 0,
		},
		throttles:     models.Throttles{0: 0, 1: 0, 2: 0, 3: 0},
		emergencyStop: false,
	}
}

func (pidcontrols *pidControls) SetFlightCommands(flightCommands models.FlightCommands) {
	pidcontrols.targetState = pidcontrols.flightControlCommandToPIDCommand(flightCommands)
	showStates(pidcontrols.targetState)
	pidcontrols.calcThrottles()
}

func (pidcontrols *pidControls) SetRotations(rotations models.ImuRotations) {
	pidcontrols.calcThrottles()
}

func (pidcontrols *pidControls) calcThrottles() {
	t := pidcontrols.targetState.throttle
	pidcontrols.throttles = models.Throttles{
		0: t,
		1: t,
		2: t,
		3: t,
	}
}

func (pidcontrols *pidControls) Throttles() models.Throttles {
	return pidcontrols.throttles
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
	return float64(joystickDigitalValue) / pidcontrols.maxJoystickDigitalValue * pidcontrols.throttleLimit
}

func (pidcontrols *pidControls) flightControlCommandToPIDCommand(c models.FlightCommands) pidState {

	return pidState{
		roll:     pidcontrols.joystickToPidValue(c.Roll, pidcontrols.roll.limit),
		pitch:    pidcontrols.joystickToPidValue(c.Pitch, pidcontrols.pitch.limit),
		yaw:      pidcontrols.joystickToPidValue(c.Yaw, pidcontrols.yaw.limit),
		throttle: pidcontrols.throttleToPidThrottle(c.Throttle),
	}
}

var lastPrint time.Time = time.Now()

func showStates(s pidState) {
	if time.Since(lastPrint) > time.Second/2 {
		lastPrint = time.Now()
		log.Printf("roll: %6.2f, pitch: %6.2f, yaw: %6.2f, throttle: %6.2f\n", s.roll, s.pitch, s.yaw, s.throttle)
	}
}

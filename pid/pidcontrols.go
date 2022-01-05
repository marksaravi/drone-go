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

func (c *pidControls) SetFlightCommands(flightCommands models.FlightCommands) {
	c.targetState = c.flightControlCommandToPIDCommand(flightCommands)
	showStates(c.state, c.targetState)
	c.calcThrottles()
}

func (c *pidControls) SetRotations(rotations models.ImuRotations) {
	c.state = pidState{
		roll:     rotations.Rotations.Roll,
		pitch:    rotations.Rotations.Pitch,
		yaw:      rotations.Rotations.Yaw,
		throttle: 0,
	}
	c.calcThrottles()
}

func (c *pidControls) calcThrottles() {
	t := c.targetState.throttle
	c.throttles = models.Throttles{
		0: t,
		1: t,
		2: t,
		3: t,
	}
}

func (c *pidControls) Throttles() models.Throttles {
	return c.throttles
}

func (c *pidControls) SetEmergencyStop(stop bool) {
	c.emergencyStop = stop
}

// func (c *pidControls) applyEmergencyStop() {
// 	c.targetState = pidTargetState{
// 		roll:     0,
// 		pitch:    0,
// 		yaw:      0,
// 		throttle: 0,
// 	}
// }

func (c *pidControls) joystickToPidValue(joystickDigitalValue uint16, maxValue float64) float64 {
	normalizedDigitalValue := float64(joystickDigitalValue) - c.maxJoystickDigitalValue/2
	return normalizedDigitalValue / c.maxJoystickDigitalValue * maxValue
}

func (c *pidControls) throttleToPidThrottle(joystickDigitalValue uint16) float64 {
	return float64(joystickDigitalValue) / c.maxJoystickDigitalValue * c.throttleLimit
}

func (c *pidControls) flightControlCommandToPIDCommand(fc models.FlightCommands) pidState {

	return pidState{
		roll:     c.joystickToPidValue(fc.Roll, c.roll.limit),
		pitch:    c.joystickToPidValue(fc.Pitch, c.pitch.limit),
		yaw:      c.joystickToPidValue(fc.Yaw, c.yaw.limit),
		throttle: c.throttleToPidThrottle(fc.Throttle),
	}
}

var lastPrint time.Time = time.Now()

func showStates(a, t pidState) {
	if time.Since(lastPrint) > time.Second/2 {
		lastPrint = time.Now()
		log.Printf("actual roll: %6.2f, pitch: %6.2f, yaw: %6.2f, throttle: %6.2f,  target roll: %6.2f, pitch: %6.2f, yaw: %6.2f, throttle: %6.2f\n    ", a.roll, a.pitch, a.yaw, a.throttle, t.roll, t.pitch, t.yaw, t.throttle)
	}
}

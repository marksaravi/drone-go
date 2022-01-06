package pid

import (
	"log"
	"math"
	"time"

	"github.com/marksaravi/drone-go/models"
)

const EMERGENCY_STOP_DURATION = time.Second * 2

type pidState struct {
	roll     float64
	pitch    float64
	yaw      float64
	throttle float64
}

type gains struct {
	P float64
	I float64
	D float64
}

type pidControls struct {
	gains                   gains
	roll                    *axisControl
	pitch                   *axisControl
	yaw                     *axisControl
	targetState             pidState
	state                   pidState
	throttles               models.Throttles
	maxJoystickDigitalValue float64
	throttleLimit           float64
	axisAlignmentAngle      float64
	calibrationGain         string
	calibrationStep         float64
	calibrationStepApplied  bool
	emergencyStopTimeout    time.Time
	emergencyStopStart      float64
}

func NewPIDControls(
	pGain, iGain, dGain float64,
	limitRoll, limitPitch, limitYaw, throttleLimit float64,
	maxJoystickDigitalValue uint16,
	axisAlignmentAngle float64,
	calibrationGain string,
	calibrationStep float64,
) *pidControls {

	return &pidControls{
		gains: gains{
			P: pGain,
			I: iGain,
			D: dGain,
		},
		roll:                    NewPIDControl(limitRoll),
		pitch:                   NewPIDControl(limitPitch),
		yaw:                     NewPIDControl(limitYaw),
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
		throttles: models.Throttles{
			BaseThrottle: 0,
			DThrottles: map[int]float32{
				0: 0,
				1: 0,
				2: 0,
				3: 0,
			},
		},
		calibrationGain:        calibrationGain,
		calibrationStep:        calibrationStep,
		calibrationStepApplied: false,
		emergencyStopTimeout:   time.Now().Add(time.Second * 86400),
		emergencyStopStart:     0,
	}
}

func (c *pidControls) SetFlightCommands(flightCommands models.FlightCommands) {
	c.targetState = c.flightControlCommandToPIDCommand(flightCommands)
	showStates(c.state, c.targetState)
	if c.calibrationGain != "none" {
		c.calibrateGain(c.calibrationGain, flightCommands.ButtonTopLeft, flightCommands.ButtonTopRight)
	}
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
	c.applyEmergencyStop()
	roll := c.roll.calc(c.state.roll, c.targetState.roll, &c.gains)
	pitch := c.pitch.calc(c.state.pitch, c.targetState.pitch, &c.gains)
	a := math.Pi / 4
	np := math.Cos(a)*roll - math.Sin(a)*pitch
	nr := math.Sin(a)*roll + math.Cos(a)*pitch
	// c.yaw.calc(0, 0, &c.gains) // not applying yaw yet
	c.throttles = models.Throttles{
		BaseThrottle: float32(c.targetState.throttle),
		DThrottles: map[int]float32{
			0: float32(-nr / 2),
			1: float32(np / 2),
			2: float32(nr / 2),
			3: float32(-np / 2),
		},
	}
}

func (c *pidControls) Throttles() models.Throttles {
	return c.throttles
}

func (c *pidControls) InitiateEmergencyStop(stop bool) {
	if stop {
		c.emergencyStopTimeout = time.Now()
		c.emergencyStopStart = c.targetState.throttle
	} else {
		c.emergencyStopTimeout = time.Now().Add(time.Second * 86400)
	}
}

func (c *pidControls) applyEmergencyStop() {
	dur := time.Since(c.emergencyStopTimeout)
	if dur > EMERGENCY_STOP_DURATION {
		dur = EMERGENCY_STOP_DURATION
	}
	if dur > 0 {
		k := float64(EMERGENCY_STOP_DURATION-dur) / float64(EMERGENCY_STOP_DURATION)

		c.targetState.throttle = c.emergencyStopStart * k
	}
}

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

func (c *pidControls) calibrateGain(gain string, down, up bool) {
	if !down && !up {
		c.calibrationStepApplied = false
		return
	}
	if c.calibrationStepApplied {
		return
	}
	var step float64 = c.calibrationStep
	if down {
		step = -step
	}
	var value float64 = 0
	switch gain {
	case "P":
		c.gains.P += step
		value = c.gains.P
	case "I":
		c.gains.I += step
		value = c.gains.I
	case "D":
		c.gains.D += step
		value = c.gains.D
	}
	log.Printf("%s Gain is changed to %6.2f\n", gain, value)
	c.calibrationStepApplied = true
}

var lastPrint time.Time = time.Now()

func showStates(a, t pidState) {
	if time.Since(lastPrint) > time.Second/2 {
		lastPrint = time.Now()
		log.Printf("actual roll: %6.2f, pitch: %6.2f, yaw: %6.2f, throttle: %6.2f,  target roll: %6.2f, pitch: %6.2f, yaw: %6.2f, throttle: %6.2f\n    ", a.roll, a.pitch, a.yaw, a.throttle, t.roll, t.pitch, t.yaw, t.throttle)
	}
}

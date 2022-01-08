package pid

import (
	"log"
	"math"
	"strings"
	"time"

	"github.com/marksaravi/drone-go/models"
)

const EMERGENCY_STOP_DURATION = time.Second * 2

type pidState struct {
	roll     float64
	pitch    float64
	yaw      float64
	throttle float64
	dt       time.Duration
}

type gains struct {
	P float64
	I float64
	D float64
}

type pidControls struct {
	gains                   gains
	yawGains                gains
	roll                    *axisControl
	pitch                   *axisControl
	yaw                     *axisControl
	targetState             pidState
	state                   pidState
	throttles               models.Throttles
	maxJoystickDigitalValue float64
	throttleLimit           float64
	safeStartThrottle       float64
	axisAlignmentAngle      float64
	calibrationGain         string
	calibrationStep         float64
	calibrationStepApplied  bool
	emergencyStopTimeout    time.Time
	emergencyStopStart      float64
}
type PIDSettings struct {
	RollPitchPGain          float64
	RollPitchIGain          float64
	RollPitchDGain          float64
	YawPGain                float64
	YawIGain                float64
	YawDGain                float64
	LimitRoll               float64
	LimitPitch              float64
	LimitYaw                float64
	LimitI                  float64
	ThrottleLimit           float64
	SafeStartThrottle       float64
	MaxJoystickDigitalValue uint16
	AxisAlignmentAngle      float64
	CalibrationGain         string
	CalibrationStep         float64
}

func NewPIDControls(settings PIDSettings) *pidControls {

	return &pidControls{
		gains: gains{
			P: settings.RollPitchPGain,
			I: settings.RollPitchIGain,
			D: settings.RollPitchDGain,
		},
		yawGains: gains{
			P: settings.YawPGain,
			I: settings.YawIGain,
			D: settings.YawDGain,
		},
		roll:                    NewPIDControl(settings.LimitRoll, settings.LimitI),
		pitch:                   NewPIDControl(settings.LimitPitch, settings.LimitI),
		yaw:                     NewPIDControl(settings.LimitYaw, settings.LimitI),
		throttleLimit:           settings.ThrottleLimit,
		safeStartThrottle:       settings.SafeStartThrottle,
		maxJoystickDigitalValue: float64(settings.MaxJoystickDigitalValue),
		axisAlignmentAngle:      settings.AxisAlignmentAngle,
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
			dt:       0,
		},
		throttles: models.Throttles{
			Throttle: 0,
			ControlVariables: map[int]float64{
				0: 0,
				1: 0,
				2: 0,
				3: 0,
			},
		},
		calibrationGain:        settings.CalibrationGain,
		calibrationStep:        settings.CalibrationStep,
		calibrationStepApplied: false,
		emergencyStopTimeout:   time.Now().Add(time.Second * 86400),
		emergencyStopStart:     0,
	}
}

func (c *pidControls) SetFlightCommands(flightCommands models.FlightCommands) {
	if c.calibrationGain != "none" {
		c.calibrateGain(c.calibrationGain, flightCommands.ButtonTopLeft, flightCommands.ButtonTopRight)
	}
	c.targetState = c.flightControlCommandToPIDCommand(flightCommands)
	showStates(c.state, c.targetState)
}

func (c *pidControls) SetRotations(rotations models.ImuRotations) {
	c.state = pidState{
		roll:     rotations.Rotations.Roll,
		pitch:    rotations.Rotations.Pitch,
		yaw:      rotations.Rotations.Yaw,
		throttle: 0,
		dt:       rotations.ReadInterval,
	}
	c.calcThrottles()
}

func (c *pidControls) calcPID(roll, pitch, yaw float64) (float64, float64, float64) {
	if c.targetState.throttle < c.safeStartThrottle {
		return 0, 0, 0
	}
	rollPID := c.roll.calc(roll, c.state.dt, &c.gains)
	pitchPID := c.pitch.calc(pitch, c.state.dt, &c.gains)
	yawPID := c.yaw.calc(c.state.yaw-c.targetState.yaw, c.state.dt, &c.yawGains)
	return rollPID, pitchPID, yawPID
}

func applySensoreZaxisRotation(roll, pitch, angle float64) (float64, float64) {
	arad := angle / 180.0 * math.Pi
	nPitch := math.Cos(arad)*roll - math.Sin(arad)*pitch
	nRoll := math.Sin(arad)*roll + math.Cos(arad)*pitch
	return nRoll, nPitch
}

func (c *pidControls) calcThrottles() {
	c.applyEmergencyStop()
	nRoll, nPitch := applySensoreZaxisRotation(c.state.roll-c.targetState.roll, c.state.pitch-c.targetState.pitch, c.axisAlignmentAngle)
	rollPID, pitchPID, _ := c.calcPID(nRoll, nPitch, c.state.yaw-c.targetState.yaw)

	c.throttles = models.Throttles{
		Throttle: c.targetState.throttle,
		ControlVariables: map[int]float64{
			0: rollPID / 2,
			1: -pitchPID / 2,
			2: -rollPID / 2,
			3: pitchPID / 2,
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
		roll:     c.joystickToPidValue(fc.Roll, c.roll.inputLimit),
		pitch:    c.joystickToPidValue(fc.Pitch, c.pitch.inputLimit),
		yaw:      c.joystickToPidValue(fc.Yaw, c.yaw.inputLimit),
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
	switch strings.ToLower(gain) {
	case "roll-p":
		c.gains.P += step
		value = c.gains.P
	case "roll-i":
		c.gains.I += step
		value = c.gains.I
	case "roll-d":
		c.gains.D += step
		value = c.gains.D
	case "yaw-p":
		c.yawGains.P += step
		value = c.yawGains.P
	case "yaw-i":
		c.yawGains.I += step
		value = c.yawGains.I
	case "yaw-d":
		c.yawGains.D += step
		value = c.yawGains.D
	}
	log.Printf("%s Gain is changed to %8.6f\n", gain, value)
	c.calibrationStepApplied = true
}

func (c *pidControls) PrintGains() {
	log.Printf("P: %8.6f, I: %8.6f, D: %8.6f, yP: %8.6f, yI: %8.6f, yD: %8.6f\n", c.gains.P, c.gains.I, c.gains.D, c.yawGains.P, c.yawGains.I, c.yawGains.D)
}

var lastPrint time.Time = time.Now()

func showStates(a, t pidState) {
	if time.Since(lastPrint) > time.Second*2 {
		lastPrint = time.Now()
		log.Printf("actual roll: %6.2f, pitch: %6.2f, yaw: %6.2f, throttle: %6.2f,  target roll: %6.2f, pitch: %6.2f, yaw: %6.2f, throttle: %6.2f\n    ", a.roll, a.pitch, a.yaw, a.throttle, t.roll, t.pitch, t.yaw, t.throttle)
	}
}

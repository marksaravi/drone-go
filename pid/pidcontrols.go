package pid

import (
	"log"
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

type PIDSettings struct {
	InputLimit float64
	PGain      float64
	IGain      float64
	DGain      float64
	ILimit     float64
}

type CalibrationSettings struct {
	Calibrating bool
	PStep       float64
	IStep       float64
	DStep       float64
}
type PIDControlSettings struct {
	Roll        PIDSettings
	Pitch       PIDSettings
	Yaw         PIDSettings
	Calibration CalibrationSettings
}
type pidControls struct {
	settings             PIDControlSettings
	roll                 *axisControl
	pitch                *axisControl
	yaw                  *axisControl
	targetState          pidState
	state                pidState
	throttles            models.Throttles
	emergencyStopTimeout time.Time
	emergencyStopStart   float64
}

func NewPIDControls(settings PIDControlSettings) *pidControls {

	return &pidControls{
		settings: settings,
		roll:     NewPIDControl(settings.Roll),
		pitch:    NewPIDControl(settings.Pitch),
		yaw:      NewPIDControl(settings.Yaw),
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
			Active: true,
			Throttles: map[int]float64{
				0: 0,
				1: 0,
				2: 0,
				3: 0,
			},
		},
		emergencyStopTimeout: time.Now().Add(time.Second * 86400),
		emergencyStopStart:   0,
	}
}

func (c *pidControls) SetFlightCommands(flightCommands models.FlightCommands) {
	// if c.calibrationGain != "none" {
	// 	c.calibrateGain(c.calibrationGain, flightCommands.ButtonTopLeft, flightCommands.ButtonTopRight)
	// }
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
	rollPID := c.roll.calc(roll, c.state.dt)
	pitchPID := c.pitch.calc(pitch, c.state.dt)
	yawPID := c.yaw.calc(c.state.yaw-c.targetState.yaw, c.state.dt)
	return rollPID, pitchPID, yawPID
}

func (c *pidControls) calcThrottles() {
	// c.applyEmergencyStop()
	// rollPID, pitchPID, yawPID := c.calcPID(
	// 	c.state.roll-c.targetState.roll,
	// 	c.state.pitch-c.targetState.pitch,
	// 	c.state.yaw-c.targetState.yaw,
	// )

	// motor0roll := rollPID / 2
	// motor3roll := rollPID / 2
	// motor1roll := -rollPID / 2
	// motor2roll := -rollPID / 2

	// motor0pitch := pitchPID / 2
	// motor1pitch := pitchPID / 2
	// motor2pitch := -pitchPID / 2
	// motor3pitch := -pitchPID / 2

	c.throttles = models.Throttles{
		Active: true,
		Throttles: map[int]float64{
			// 0: motor0roll + motor0pitch + yawPID/2,
			// 1: motor1roll + motor1pitch - yawPID/2,
			// 2: motor2roll + motor2pitch + yawPID/2,
			// 3: motor3roll + motor3pitch - yawPID/2,
			0: 0,
			1: 0,
			2: 0,
			3: 0,
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

func (c *pidControls) flightControlCommandToPIDCommand(fc models.FlightCommands) pidState {
	return pidState{
		roll:     0, //c.joystickToPidValue(fc.Roll, c.roll.inputLimit),
		pitch:    0, //c.joystickToPidValue(fc.Pitch, c.pitch.inputLimit),
		yaw:      0, //c.joystickToPidValue(fc.Yaw, c.yaw.inputLimit),
		throttle: 0, //c.throttleToPidThrottle(fc.Throttle),
	}
}

func (c *pidControls) calibrateGain(gain string, down, up bool) {
	// addStep := func(x, step float64) float64 {
	// 	nvalue := x + step
	// 	if nvalue < 0 {
	// 		nvalue = x
	// 	}
	// 	log.Printf("%s Gain is changed to %8.6f\n", gain, nvalue)
	// 	return nvalue
	// }
	// if !down && !up {
	// 	c.calibrationStepApplied = false
	// 	return
	// }
	// if c.calibrationStepApplied {
	// 	return
	// }
	// var step float64 = c.calibrationStep
	// if down {
	// 	step = -step
	// }

	// switch strings.ToLower(gain) {
	// case "roll-p":
	// 	c.gains.P = addStep(c.gains.P, step)
	// case "roll-i":
	// 	c.gains.I = addStep(c.gains.I, step)
	// case "roll-d":
	// 	c.gains.D = addStep(c.gains.D, step)
	// case "yaw-p":
	// 	c.yawGains.P = addStep(c.yawGains.P, step)
	// case "yaw-i":
	// 	c.yawGains.I = addStep(c.yawGains.I, step)
	// case "yaw-d":
	// 	c.yawGains.D = addStep(c.yawGains.D, step)
	// }
	// c.calibrationStepApplied = true
}

func (c *pidControls) PrintGains() {
	// log.Printf("P: %8.6f, I: %8.6f, D: %8.6f, yP: %8.6f, yI: %8.6f, yD: %8.6f\n", c.gains.P, c.gains.I, c.gains.D, c.yawGains.P, c.yawGains.I, c.yawGains.D)
}

var lastPrint time.Time = time.Now()

func showStates(a, t pidState) {
	if time.Since(lastPrint) > time.Second*2 {
		lastPrint = time.Now()
		log.Printf("actual roll: %6.2f, pitch: %6.2f, yaw: %6.2f, throttle: %6.2f,  target roll: %6.2f, pitch: %6.2f, yaw: %6.2f, throttle: %6.2f\n    ", a.roll, a.pitch, a.yaw, a.throttle, t.roll, t.pitch, t.yaw, t.throttle)
	}
}

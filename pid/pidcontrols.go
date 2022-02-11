package pid

import (
	"log"
	"time"

	"github.com/marksaravi/drone-go/models"
)

const EMERGENCY_STOP_DURATION = time.Second * 2

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
	settings    PIDControlSettings
	roll        *axisControl
	pitch       *axisControl
	yaw         *axisControl
	targetState models.PIDState
	state       models.PIDState
	throttles   map[int]float64
}

func NewPIDControls(settings PIDControlSettings) *pidControls {

	return &pidControls{
		settings: settings,
		roll:     NewPIDControl(settings.Roll),
		pitch:    NewPIDControl(settings.Pitch),
		yaw:      NewPIDControl(settings.Yaw),
		targetState: models.PIDState{
			Roll:     0,
			Pitch:    0,
			Yaw:      0,
			Throttle: 0,
		},
		state: models.PIDState{
			Roll:     0,
			Pitch:    0,
			Yaw:      0,
			Throttle: 0,
		},
		throttles: map[int]float64{
			0: 0,
			1: 0,
			2: 0,
			3: 0,
		},
	}
}

func (c *pidControls) SetPIDTargetState(state models.PIDState) {
	c.targetState = state
	showStates(c.state, c.targetState)
}

func (c *pidControls) SetRotations(rotations models.ImuRotations) {
	c.state = models.PIDState{
		Roll:     rotations.Rotations.Roll,
		Pitch:    rotations.Rotations.Pitch,
		Yaw:      rotations.Rotations.Yaw,
		Throttle: 0,
		Dt:       rotations.ReadInterval,
	}
	c.calcThrottles()
}

func (c *pidControls) calcPID(roll, pitch, yaw float64, dt time.Duration) (float64, float64, float64) {
	rollPID := c.roll.calc(roll, dt)
	pitchPID := c.pitch.calc(pitch, dt)
	yawPID := c.yaw.calc(yaw, dt)
	return rollPID, pitchPID, yawPID
}

func (c *pidControls) calcThrottles() {
	c.calcPID(
		c.state.Roll-c.targetState.Roll,
		c.state.Pitch-c.targetState.Pitch,
		c.state.Yaw-c.targetState.Yaw,
		c.state.Dt,
	)

	c.throttles = map[int]float64{
		0: 0,
		1: 0,
		2: 0,
		3: 0,
	}
}

func (c *pidControls) Throttles() map[int]float64 {
	return c.throttles
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

func showStates(a, t models.PIDState) {
	if time.Since(lastPrint) > time.Second*2 {
		lastPrint = time.Now()
		// log.Printf("actual roll: %6.2f, pitch: %6.2f, yaw: %6.2f, throttle: %6.2f,  target roll: %6.2f, pitch: %6.2f, yaw: %6.2f, throttle: %6.2f\n    ", a.Roll, a.Pitch, a.Yaw, a.Throttle, t.Roll, t.Pitch, t.Yaw, t.Throttle)
		log.Printf("throttle: %6.2f,  target roll: %6.2f, pitch: %6.2f, yaw: %6.2f\n    ", t.Throttle, t.Roll, t.Pitch, t.Yaw)
	}
}

package pid

import (
	"log"
	"time"

	"github.com/marksaravi/drone-go/models"
)

const EMERGENCY_STOP_DURATION = time.Second * 2

type PIDSettings struct {
	MaxOutputToMaxThrottleRatio float64
	PGain                       float64
	IGain                       float64
	DGain                       float64
}

type CalibrationSettings struct {
	Calibrating string
	PStep       float64
	IStep       float64
	DStep       float64
}
type pidControls struct {
	rollPIDControl  *pidControl
	pitchPIDControl *pidControl
	yawPIDControl   *pidControl
	targetStates    models.PIDState
	states          models.PIDState
	throttle        float64
	throttles       map[int]float64
	calibration     CalibrationSettings
}

func NewPIDControls(
	rollPIDSettings PIDSettings,
	pitchPIDSettings PIDSettings,
	yawPIDSettings PIDSettings,
	maxThrottle float64,
	maxItoMaxOutputRatio float64,
	calibration CalibrationSettings,
) *pidControls {
	return &pidControls{
		calibration:     calibration,
		rollPIDControl:  NewPIDControl(rollPIDSettings, maxThrottle, maxItoMaxOutputRatio),
		pitchPIDControl: NewPIDControl(pitchPIDSettings, maxThrottle, maxItoMaxOutputRatio),
		yawPIDControl:   NewPIDControl(yawPIDSettings, maxThrottle, maxItoMaxOutputRatio),
		targetStates: models.PIDState{
			Roll:         0,
			Pitch:        0,
			Yaw:          0,
			ReadTime:     time.Now(),
			ReadInterval: 0,
		},
		states: models.PIDState{
			Roll:         0,
			Pitch:        0,
			Yaw:          0,
			ReadTime:     time.Now(),
			ReadInterval: 0,
		},
		throttles: map[int]float64{
			0: 0,
			1: 0,
			2: 0,
			3: 0,
		},
	}
}

func (c *pidControls) SetTargetStates(states models.PIDState, throttle float64) {
	c.targetStates = states
	c.throttle = throttle
}

func (c *pidControls) SetStates(rotations models.ImuRotations) {
	c.states = models.PIDState{
		Roll:         rotations.Rotations.Roll,
		Pitch:        rotations.Rotations.Pitch,
		Yaw:          rotations.Rotations.Yaw,
		ReadTime:     rotations.ReadTime,
		ReadInterval: rotations.ReadInterval,
	}

}

func (c *pidControls) calcPIDs(roll, pitch, yaw float64, dt time.Duration) (float64, float64, float64) {
	rollPID := c.rollPIDControl.calc(roll, dt)
	pitchPID := c.pitchPIDControl.calc(pitch, dt)
	yawPID := c.yawPIDControl.calc(yaw, dt)
	return rollPID, pitchPID, yawPID
}

func (c *pidControls) calcThrottles() {
	c.calcPIDs(
		c.states.Roll-c.targetStates.Roll,
		c.states.Pitch-c.targetStates.Pitch,
		c.states.Yaw-c.targetStates.Yaw,
		c.states.ReadInterval,
	)

	c.throttles = map[int]float64{
		0: 0,
		1: 0,
		2: 0,
		3: 0,
	}
}

func (c *pidControls) GetThrottles() map[int]float64 {
	c.calcThrottles()
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

func showStates(a, t models.PIDState, throttle float64) {
	if time.Since(lastPrint) > time.Second*2 {
		lastPrint = time.Now()
		// log.Printf("actual roll: %6.2f, pitch: %6.2f, yaw: %6.2f, throttle: %6.2f,  target roll: %6.2f, pitch: %6.2f, yaw: %6.2f, throttle: %6.2f\n    ", a.Roll, a.Pitch, a.Yaw, a.Throttle, t.Roll, t.Pitch, t.Yaw, t.Throttle)
		log.Printf("throttle: %6.2f,  target roll: %6.2f, pitch: %6.2f, yaw: %6.2f\n    ", throttle, t.Roll, t.Pitch, t.Yaw)
	}
}

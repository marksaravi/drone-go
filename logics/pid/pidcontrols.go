package pid

import (
	"fmt"
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
	Gain        string
	PStep       float64
	IStep       float64
	DStep       float64
}
type pidControls struct {
	rollPIDControl  *pidControl
	pitchPIDControl *pidControl
	yawPIDControl   *pidControl
	targetStates    models.Rotations
	states          models.Rotations
	dt              time.Duration
	throttle        float64
	throttles       map[int]float64
	calibration     CalibrationSettings
}

func NewPIDControls(
	rollPIDSettings PIDSettings,
	pitchPIDSettings PIDSettings,
	yawPIDSettings PIDSettings,
	maxThrottle float64,
	minThrottle float64,
	maxItoMaxOutputRatio float64,
	calibration CalibrationSettings,
) *pidControls {
	return &pidControls{
		calibration:     calibration,
		rollPIDControl:  NewPIDControl("Roll", rollPIDSettings, maxThrottle, minThrottle, maxItoMaxOutputRatio),
		pitchPIDControl: NewPIDControl("Pitch", pitchPIDSettings, maxThrottle, minThrottle, maxItoMaxOutputRatio),
		yawPIDControl:   NewPIDControl("Yaw", yawPIDSettings, maxThrottle, minThrottle, maxItoMaxOutputRatio),
		targetStates: models.Rotations{
			Roll:  0,
			Pitch: 0,
			Yaw:   0,
		},
		states: models.Rotations{
			Roll:  0,
			Pitch: 0,
			Yaw:   0,
		},
		dt: 0,
		throttles: map[int]float64{
			0: 0,
			1: 0,
			2: 0,
			3: 0,
		},
	}
}

func (c *pidControls) SetTargetStates(states models.Rotations, throttle float64) {
	c.targetStates = states
	c.throttle = throttle
}

func (c *pidControls) SetStates(rotations models.Rotations, dt time.Duration) {
	c.states = rotations
	c.dt = dt

}

func (c *pidControls) calcPIDs(errRoll, errPitch, errYaw float64, dt time.Duration) (rollFeedback, pitchFeedback, yawFeedback float64) {
	rollFeedback = c.rollPIDControl.calcPIDFeedback(errRoll, dt)
	pitchFeedback = c.pitchPIDControl.calcPIDFeedback(errPitch, dt)
	yawFeedback = c.yawPIDControl.calcPIDFeedback(errYaw, dt)
	return
}

func (c *pidControls) calcThrottles() {
	// rollFeedback, pitchFeedback, yawFeedback := c.calcPIDs(
	// 	c.states.Roll-c.targetStates.Roll,
	// 	c.states.Pitch-c.targetStates.Pitch,
	// 	c.states.Yaw-c.targetStates.Yaw,
	// 	c.states.ReadInterval,
	// )

	c.throttles = map[int]float64{
		0: c.throttle, //- pitchFeedback - yawFeedback,
		1: c.throttle, //- rollFeedback + yawFeedback,
		2: c.throttle, //+ pitchFeedback - yawFeedback,
		3: c.throttle, //+ rollFeedback + yawFeedback,
	}
}

func (c *pidControls) GetThrottles(isSafeStarted bool) models.Throttles {
	if !isSafeStarted {
		return models.Throttles{
			BaseThrottle: 0,
			Throttles:    map[int]float64{0: 0, 1: 0, 2: 0, 3: 0},
		}
	}
	c.calcThrottles()
	return models.Throttles{
		BaseThrottle: c.throttle,
		Throttles:    c.throttles,
	}
}

func (c *pidControls) Calibrate(down, up bool) {
	if c.calibration.Calibrating == "none" {
		return
	}
	var pidcontrol *pidControl
	switch c.calibration.Calibrating {
	case "roll":
		pidcontrol = c.rollPIDControl
	case "pitch":
		pidcontrol = c.pitchPIDControl
	case "yaw":
		pidcontrol = c.yawPIDControl
	}
	var sign float64 = 1
	if down {
		sign = -1
	}
	switch c.calibration.Gain {
	case "p":
		pidcontrol.pGain += c.calibration.PStep * sign
	case "i":
		pidcontrol.iGain += c.calibration.IStep * sign
	case "d":
		pidcontrol.dGain += c.calibration.DStep * sign
	}
	if up || down {
		printGains(pidcontrol)
	}

}

func printGains(p *pidControl) {
	log.Printf("%8s P: %8.6f, I: %8.6f, D: %8.6f\n", p.name, p.pGain, p.iGain, p.dGain)
}

func (p *pidControls) PrintGains() {
	if p.calibration.Calibrating == "none" {
		return
	}
	fmt.Println()
	printGains(p.rollPIDControl)
	printGains(p.pitchPIDControl)
	printGains(p.yawPIDControl)
	fmt.Println()
}

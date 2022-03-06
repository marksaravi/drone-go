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
		rollPIDControl:  NewPIDControl("Roll", rollPIDSettings, maxThrottle, maxItoMaxOutputRatio),
		pitchPIDControl: NewPIDControl("Pitch", pitchPIDSettings, maxThrottle, maxItoMaxOutputRatio),
		yawPIDControl:   NewPIDControl("Yaw", yawPIDSettings, maxThrottle, maxItoMaxOutputRatio),
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

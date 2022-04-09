package pid

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/marksaravi/drone-go/models"
)

const EMERGENCY_STOP_DURATION = time.Second * 2
const HEADING_NOT_SET float64 = 1000000

type PIDSettings struct {
	PGain float64
	IGain float64
	DGain float64
	MaxI  float64
}

type CalibrationSettings struct {
	Calibrating string
	Gain        string
	PStep       float64
	IStep       float64
	DStep       float64
}

type pidControls struct {
	rollPIDControl          *pidControl
	pitchPIDControl         *pidControl
	yawPIDControl           *pidControl
	targetStates            models.RotationsAroundAxis
	heading                 float64
	arm_0_2_ThrottleEnabled bool
	arm_1_3_ThrottleEnabled bool
	minThrottle             float64
	calibration             CalibrationSettings
	calibrationApplied      bool
}

func NewPIDControls(
	rollPIDSettings PIDSettings,
	pitchPIDSettings PIDSettings,
	yawPIDSettings PIDSettings,
	arm_0_2_ThrottleEnabled bool,
	arm_1_3_ThrottleEnabled bool,
	minThrottle float64,
	calibration CalibrationSettings,
) *pidControls {
	return &pidControls{
		calibration:             calibration,
		rollPIDControl:          NewPIDControl("Roll", rollPIDSettings),
		pitchPIDControl:         NewPIDControl("Pitch", pitchPIDSettings),
		yawPIDControl:           NewPIDControl("Yaw", yawPIDSettings),
		arm_0_2_ThrottleEnabled: arm_0_2_ThrottleEnabled,
		arm_1_3_ThrottleEnabled: arm_1_3_ThrottleEnabled,
		targetStates: models.RotationsAroundAxis{
			Roll:  0,
			Pitch: 0,
			Yaw:   0,
		},
		heading:            HEADING_NOT_SET,
		minThrottle:        minThrottle,
		calibrationApplied: false,
	}
}

func (c *pidControls) SetTargetStates(states models.RotationsAroundAxis) {
	c.targetStates = states
}

func (c *pidControls) calcAxisFeedbacks(rollError, pitchError, yawError float64, dt time.Duration) (rollFeedback, pitchFeedback, yawFeedback float64) {
	rollFeedback = c.rollPIDControl.calcPIDFeedback(rollError, dt)
	pitchFeedback = c.pitchPIDControl.calcPIDFeedback(pitchError, dt)
	yawFeedback = c.yawPIDControl.calcPIDFeedback(yawError, dt)
	return
}

func (c *pidControls) calcArmsFeedbacks(rollFeedback, pitchFeedback, yawFeedback float64) (arms [4]float64) {
	arms = [4]float64{0, 0, 0, 0}
	arms[0] = +pitchFeedback - yawFeedback
	arms[1] = +rollFeedback + yawFeedback
	arms[2] = -pitchFeedback - yawFeedback
	arms[3] = -rollFeedback + yawFeedback
	return
}

func (c *pidControls) calculateThrottles(throttle float64, armsFeedback [4]float64) models.Throttles {
	applyFeedbacks := float64(1)
	if throttle < c.minThrottle {
		applyFeedbacks = 0
		c.reset()
	}
	arm_0_2_enabled := float64(0)
	if c.arm_0_2_ThrottleEnabled {
		arm_0_2_enabled = 1.0
	}
	arm_1_3_enabled := float64(0)
	if c.arm_1_3_ThrottleEnabled {
		arm_1_3_enabled = 1.0
	}
	return models.Throttles{
		BaseThrottle: throttle,
		Throttles: map[int]float64{
			0: (throttle + armsFeedback[0]*applyFeedbacks) * arm_0_2_enabled,
			1: (throttle + armsFeedback[1]*applyFeedbacks) * arm_1_3_enabled,
			2: (throttle + armsFeedback[2]*applyFeedbacks) * arm_0_2_enabled,
			3: (throttle + armsFeedback[3]*applyFeedbacks) * arm_1_3_enabled,
		},
	}

}

func (c *pidControls) GetThrottles(throttle float64, rotations models.RotationsAroundAxis, dt time.Duration) models.Throttles {
	if c.heading == HEADING_NOT_SET {
		c.heading = rotations.Yaw
	}
	rollError := c.targetStates.Roll - rotations.Roll
	pitchError := c.targetStates.Pitch - rotations.Pitch
	yawError := c.heading - rotations.Yaw
	if math.Abs(c.targetStates.Yaw) > 1 {
		yawError = c.targetStates.Yaw
		c.heading = rotations.Yaw
	}
	rollFeedback, pitchFeedback, yawFeedback := c.calcAxisFeedbacks(rollError, pitchError, yawError, dt)
	// utils.PrintIntervally(fmt.Sprintf("errors roll: %7.3f pitch: %7.3f yaw: %7.3f\n", rollError, pitchError, yawError), "yawfeedback", time.Second/2, false)
	armsFeedback := c.calcArmsFeedbacks(rollFeedback, pitchFeedback, yawFeedback)

	throttles := c.calculateThrottles(throttle, armsFeedback)
	return throttles
}

func (c *pidControls) reset() {
	c.rollPIDControl.reset()
	c.pitchPIDControl.reset()
	c.yawPIDControl.reset()
}
func (c *pidControls) setCalibrationGain(axis, gain string, up, down bool) {
	var pidcontrol *pidControl
	switch axis {
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
	switch gain {
	case "p":
		pidcontrol.pGain += c.calibration.PStep * sign
	case "i":
		pidcontrol.iGain += c.calibration.IStep * sign
	case "d":
		pidcontrol.dGain += c.calibration.DStep * sign
	}
	c.calibrationApplied = true
	if up || down {
		fmt.Println(c.calibration)
		printGains(pidcontrol)
	}
}

func (c *pidControls) Calibrate(up, down bool) {
	if !down && !up {
		c.calibrationApplied = false
		return
	}
	if c.calibrationApplied {
		return
	}
	if c.calibration.Calibrating == "none" {
		return
	}

	switch c.calibration.Calibrating {
	case "roll-pitch":
		c.setCalibrationGain("roll", c.calibration.Gain, up, down)
		c.setCalibrationGain("pitch", c.calibration.Gain, up, down)
	case "roll":
		c.setCalibrationGain("roll", c.calibration.Gain, up, down)
	case "pitch":
		c.setCalibrationGain("pitch", c.calibration.Gain, up, down)
	case "yaw":
		c.setCalibrationGain("yaw", c.calibration.Gain, up, down)
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

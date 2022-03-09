package pid

import (
	"fmt"
	"log"
	"time"

	"github.com/marksaravi/drone-go/models"
)

const EMERGENCY_STOP_DURATION = time.Second * 2

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
	rollPIDControl        *pidControl
	pitchPIDControl       *pidControl
	yawPIDControl         *pidControl
	targetStates          models.Rotations
	arm_0_2_ThrottleRatio float64
	arm_1_3_ThrottleRatio float64
	minThrottle           float64
	calibration           CalibrationSettings
	calibrationApplied    bool
}

func NewPIDControls(
	rollPIDSettings PIDSettings,
	pitchPIDSettings PIDSettings,
	yawPIDSettings PIDSettings,
	arm_0_2_ThrottleRatio float64,
	arm_1_3_ThrottleRatio float64,
	minThrottle float64,
	calibration CalibrationSettings,
) *pidControls {
	return &pidControls{
		calibration:           calibration,
		rollPIDControl:        NewPIDControl("Roll", rollPIDSettings),
		pitchPIDControl:       NewPIDControl("Pitch", pitchPIDSettings),
		yawPIDControl:         NewPIDControl("Yaw", yawPIDSettings),
		arm_0_2_ThrottleRatio: arm_0_2_ThrottleRatio,
		arm_1_3_ThrottleRatio: arm_1_3_ThrottleRatio,
		targetStates: models.Rotations{
			Roll:  0,
			Pitch: 0,
			Yaw:   0,
		},
		minThrottle:        minThrottle,
		calibrationApplied: false,
	}
}

func (c *pidControls) SetTargetStates(states models.Rotations) {
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
	apply := float64(1)
	if throttle < c.minThrottle {
		apply = 0
		c.reset()
	}
	return models.Throttles{
		BaseThrottle: throttle,
		Throttles: map[int]float64{
			0: throttle*c.arm_0_2_ThrottleRatio + armsFeedback[0]*apply,
			1: throttle*c.arm_1_3_ThrottleRatio + armsFeedback[1]*apply,
			2: throttle*c.arm_0_2_ThrottleRatio + armsFeedback[2]*apply,
			3: throttle*c.arm_1_3_ThrottleRatio + armsFeedback[3]*apply,
		},
	}

}

var ts time.Time = time.Now()

func (c *pidControls) GetThrottles(throttle float64, rotations models.Rotations, dt time.Duration) models.Throttles {
	rollError := c.targetStates.Roll - rotations.Roll
	pitchError := c.targetStates.Pitch - rotations.Pitch
	yawError := c.targetStates.Yaw - rotations.Yaw
	rollFeedback, pitchFeedback, yawFeedback := c.calcAxisFeedbacks(rollError, pitchError, yawError, dt)
	armsFeedback := c.calcArmsFeedbacks(rollFeedback, pitchFeedback, yawFeedback)

	throttles := c.calculateThrottles(throttle, armsFeedback)
	if time.Since(ts) >= time.Second/2 {
		ts = time.Now()
		log.Printf(" pid: { roll: %8.4f, pitch: %8.4f, yaw: %8.4f } throttles: {0: %8.4f, 1: %8.4f, 2: %8.4f, 3: %8.4f}\n", rollFeedback, pitchFeedback, yawFeedback, throttles.Throttles[0], throttles.Throttles[1], throttles.Throttles[2], throttles.Throttles[3])
	}
	return throttles
}

func (c *pidControls) reset() {
	c.rollPIDControl.reset()
	c.pitchPIDControl.reset()
	c.yawPIDControl.reset()
}

func (c *pidControls) Calibrate(down, up bool) {
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
	c.calibrationApplied = true
	if up || down {
		fmt.Println(c.calibration)
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

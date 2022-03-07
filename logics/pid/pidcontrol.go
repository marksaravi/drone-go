package pid

import (
	"time"

	"github.com/marksaravi/drone-go/utils"
)

type pidControl struct {
	name          string
	pGain         float64
	iGain         float64
	dGain         float64
	maxOutput     float64
	minOutput     float64
	maxI          float64
	previousInput float64
	iMemory       float64
}

func NewPIDControl(name string, settings PIDSettings, maxThrottle float64, minThrottle float64, maxIToMaxOutputRatio float64) *pidControl {
	maxPIDOutput := settings.MaxOutputToMaxThrottleRatio * maxThrottle
	return &pidControl{
		name:          name,
		pGain:         settings.PGain,
		iGain:         settings.IGain,
		dGain:         settings.DGain,
		maxOutput:     maxPIDOutput,
		minOutput:     minThrottle,
		maxI:          maxPIDOutput * maxIToMaxOutputRatio,
		previousInput: 0,
		iMemory:       0,
	}
}

func (pidcontrol *pidControl) reset() {
	pidcontrol.iMemory = 0
}

func (pidcontrol *pidControl) getP(input float64) float64 {
	return input * pidcontrol.pGain
}

func (pidcontrol *pidControl) getI(input float64, dt time.Duration) float64 {
	pidcontrol.iMemory = utils.ApplyLimits(input*pidcontrol.iGain*float64(dt)/1000000000+pidcontrol.iMemory, -pidcontrol.maxI, pidcontrol.maxI)
	return pidcontrol.iMemory
}

func (pidcontrol *pidControl) getD(input float64, dt time.Duration) float64 {
	d := (input - pidcontrol.previousInput) / float64(dt) * 1000000000 * pidcontrol.dGain

	pidcontrol.previousInput = input
	return d
}

func (pidcontrol *pidControl) calcPIDFeedback(input float64, dt time.Duration) float64 {
	p := pidcontrol.getP(input)
	i := pidcontrol.getI(input, dt)
	d := pidcontrol.getD(input, dt)
	sum := utils.ApplyLimits(p+i+d, pidcontrol.minOutput, pidcontrol.maxOutput)
	return sum
}

package pid

import (
	"time"
)

type pidControl struct {
	name          string
	pGain         float64
	iGain         float64
	dGain         float64
	maxOutput     float64
	maxI          float64
	previousInput float64
	iMemory       float64
}

func NewPIDControl(name string, settings PIDSettings, maxThrottle float64, maxItoMaxOutputRatio float64) *pidControl {
	maxOutput := settings.MaxOutputToMaxThrottleRatio * maxThrottle
	return &pidControl{
		name:          name,
		pGain:         settings.PGain,
		iGain:         settings.IGain,
		dGain:         settings.DGain,
		maxOutput:     maxOutput,
		maxI:          maxOutput * maxItoMaxOutputRatio,
		previousInput: 0,
		iMemory:       0,
	}
}

func (pidcontrol *pidControl) getP(input float64) float64 {
	return input * pidcontrol.pGain
}

func (pidcontrol *pidControl) getI(input float64, dt time.Duration) float64 {
	pidcontrol.iMemory = limitToMax(input*pidcontrol.iGain*float64(dt)/1000000000+pidcontrol.iMemory, pidcontrol.maxI)

	return pidcontrol.iMemory
}

func (pidcontrol *pidControl) getD(input float64, dt time.Duration) float64 {
	d := (input - pidcontrol.previousInput) / float64(dt) * 1000000000 * pidcontrol.dGain

	pidcontrol.previousInput = input
	return d
}

func (pidcontrol *pidControl) calc(input float64, dt time.Duration) float64 {
	p := pidcontrol.getP(input)
	i := pidcontrol.getI(input, dt)
	d := pidcontrol.getD(input, dt)
	sum := p + i + d
	return sum
}

func limitToMax(x, limit float64) float64 {
	if x > limit {
		return limit
	}
	if x < -limit {
		return -limit
	}
	return x
}

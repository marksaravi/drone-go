package pid

import (
	"time"
)

type pidControl struct {
	name       string
	pGain      float64
	iGain      float64
	dGain      float64
	maxI       float64
	dPrevError float64
	iMemory    float64
}

func NewPIDControl(name string, settings PIDSettings) *pidControl {
	return &pidControl{
		name:       name,
		pGain:      settings.PGain,
		iGain:      settings.IGain,
		dGain:      settings.DGain,
		maxI:       settings.MaxI,
		dPrevError: 0,
		iMemory:    0,
	}
}

func (pidcontrol *pidControl) reset() {
	pidcontrol.iMemory = 0
}

func (pidcontrol *pidControl) getP(inputError float64) float64 {
	return inputError * pidcontrol.pGain
}

func applyLimits(x, min, max float64) float64 {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

func (pidcontrol *pidControl) getI(inputError float64, dt time.Duration) float64 {
	pidcontrol.iMemory = applyLimits(inputError*pidcontrol.iGain*float64(dt)/1000000000+pidcontrol.iMemory, -pidcontrol.maxI, pidcontrol.maxI)
	return pidcontrol.iMemory
}

func (pidcontrol *pidControl) getD(inputError float64, dt time.Duration) float64 {
	d := (inputError - pidcontrol.dPrevError) / float64(dt) * 1000000000 * pidcontrol.dGain

	pidcontrol.dPrevError = inputError
	return d
}

func (pidcontrol *pidControl) calcPIDFeedback(inputError float64, dt time.Duration) float64 {
	p := pidcontrol.getP(inputError)
	i := pidcontrol.getI(inputError, dt)
	d := pidcontrol.getD(inputError, dt)

	return p + i + d
}

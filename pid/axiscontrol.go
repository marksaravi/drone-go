package pid

import (
	"time"
)

type axisControl struct {
	settings      PIDSettings
	previousInput float64
	iMemory       float64
}

func NewPIDControl(settings PIDSettings) *axisControl {
	return &axisControl{
		settings:      settings,
		previousInput: 0,
		iMemory:       0,
	}
}

func (ac *axisControl) getP(input float64) float64 {
	return input * ac.settings.PGain
}

func (ac *axisControl) getI(input float64, dt time.Duration) float64 {
	ac.iMemory = limitToMax(input*ac.settings.IGain*float64(dt)/1000000000+ac.iMemory, ac.settings.ILimit)

	return ac.iMemory
}

func (ac *axisControl) getD(input float64, dt time.Duration) float64 {
	d := (input - ac.previousInput) / float64(dt) * 1000000000 * ac.settings.DGain

	ac.previousInput = input
	return d
}

func (ac *axisControl) calc(input float64, dt time.Duration) float64 {
	p := ac.getP(input)
	i := ac.getI(input, dt)
	d := ac.getD(input, dt)
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

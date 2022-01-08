package pid

import (
	"time"
)

type axisControl struct {
	previousInput float64
	inputLimit    float64
	iMemoryLimit  float64
	iMemory       float64
}

func NewPIDControl(inputLimit, limitI float64) *axisControl {
	return &axisControl{
		inputLimit:    inputLimit,
		iMemoryLimit:  limitI,
		previousInput: 0,
		iMemory:       0,
	}
}

func (ac *axisControl) getP(input, gain float64) float64 {
	return input * gain
}

func (ac *axisControl) getI(input float64, gain float64, dt time.Duration) float64 {
	ac.iMemory = limitToMax(input*gain*float64(dt)/1000000000+ac.iMemory, ac.iMemoryLimit)

	return ac.iMemory
}

func (ac *axisControl) getD(input, gain float64, dt time.Duration) float64 {
	d := (input - ac.previousInput) / float64(dt) * 1000000000 * gain

	ac.previousInput = input
	return d
}

func (ac *axisControl) calc(input float64, dt time.Duration, gains *gains) float64 {
	p := ac.getP(input, gains.P)
	i := ac.getI(input, gains.I, dt)
	d := ac.getD(input, gains.D, dt)
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

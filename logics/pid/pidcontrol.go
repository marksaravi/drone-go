package pid

import (
	"time"

	"github.com/marksaravi/drone-go/models"
)

type pidControl struct {
	settings      models.PIDSettings
	previousInput float64
	iMemory       float64
}

func NewPIDControl(settings models.PIDSettings) *pidControl {
	return &pidControl{
		settings:      settings,
		previousInput: 0,
		iMemory:       0,
	}
}

func (pidcontrol *pidControl) getP(input float64) float64 {
	return input * pidcontrol.settings.PGain
}

func (pidcontrol *pidControl) getI(input float64, dt time.Duration) float64 {
	pidcontrol.iMemory = limitToMax(input*pidcontrol.settings.IGain*float64(dt)/1000000000+pidcontrol.iMemory, pidcontrol.settings.ILimit)

	return pidcontrol.iMemory
}

func (pidcontrol *pidControl) getD(input float64, dt time.Duration) float64 {
	d := (input - pidcontrol.previousInput) / float64(dt) * 1000000000 * pidcontrol.settings.DGain

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

package pid

import (
	"log"
	"time"

	"github.com/marksaravi/drone-go/utils"
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

func NewPIDControl(name string, settings PIDSettings, maxThrottle float64, maxIToMaxOutputRatio float64) *pidControl {
	maxPIDOutput := settings.MaxOutputToMaxThrottleRatio * maxThrottle
	return &pidControl{
		name:          name,
		pGain:         settings.PGain,
		iGain:         settings.IGain,
		dGain:         settings.DGain,
		maxOutput:     maxPIDOutput,
		maxI:          maxPIDOutput * maxIToMaxOutputRatio,
		previousInput: 0,
		iMemory:       0,
	}
}

func (pidcontrol *pidControl) reset() {
	pidcontrol.iMemory = 0
}

func (pidcontrol *pidControl) getP(inputError float64) float64 {
	return inputError * pidcontrol.pGain
}

func (pidcontrol *pidControl) getI(inputError float64, dt time.Duration) float64 {
	pidcontrol.iMemory = utils.ApplyLimits(inputError*pidcontrol.iGain*float64(dt)/1000000000+pidcontrol.iMemory, -pidcontrol.maxI, pidcontrol.maxI)
	return pidcontrol.iMemory
}

func (pidcontrol *pidControl) getD(inputError float64, dt time.Duration) float64 {
	d := (inputError - pidcontrol.previousInput) / float64(dt) * 1000000000 * pidcontrol.dGain

	pidcontrol.previousInput = inputError
	return d
}

func (pidcontrol *pidControl) calcPIDFeedback(inputError float64, dt time.Duration) float64 {
	p := pidcontrol.getP(inputError)
	i := pidcontrol.getI(inputError, dt)
	d := pidcontrol.getD(inputError, dt)
	utils.PrintByInterval("pidfeedbacks", time.Second/10, func() {
		log.Printf("feedbacks { roll: %8.4f, pitch: %8.4f, yaw: %8.4f\n", p, i, d)
	})

	return p + i + d
}

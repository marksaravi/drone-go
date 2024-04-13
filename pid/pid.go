package pid

import (
	"time"

	"github.com/marksaravi/drone-go/utils"
)

const MAX_SAMPLING_RATE = 10000
const MIN_DT = time.Second / MAX_SAMPLING_RATE

type PIDSettings struct {
	PGain               float64
	IGain               float64
	DGain               float64
	MaxError            float64
	MaxIntegrationValue float64
	MinProcessVariable  float64
	MaxWeightedSum      float64
}

type PIDControl struct {
	pGain               float64
	iGain               float64
	dGain               float64
	maxError            float64
	maxIntegrationValue float64
	minProcessVariable  float64
	maxWeightedSum      float64
	setPoint            float64
	feedbackTime        time.Time
	dt                  time.Duration
	errorValue          float64
	prevErrorValue      float64
	pControlVariable    float64
	iControlVariable    float64
	dControlVariable    float64
	weightedSum         float64
}

func NewPIDControl(settings PIDSettings) *PIDControl {
	pid := &PIDControl{
		pGain:               settings.PGain,
		iGain:               settings.IGain,
		dGain:               settings.DGain,
		maxError:            settings.MaxError,
		maxIntegrationValue: settings.MaxIntegrationValue,
		minProcessVariable:  settings.MinProcessVariable,
		maxWeightedSum:      settings.MaxWeightedSum,
		setPoint:            0,
		feedbackTime:        time.Now().Add(time.Second * 32000000),
		dt:                  0,
		errorValue:          0,
		prevErrorValue:      0,
		pControlVariable:    0,
		iControlVariable:    0,
		dControlVariable:    0,
		weightedSum:         0,
	}
	return pid
}

func (pid *PIDControl) CalculateControlVariable(processVariable float64, t time.Time) float64 {
	if t.Sub(pid.feedbackTime) < 0 {
		pid.feedbackTime = t
	}
	pid.dt = t.Sub(pid.feedbackTime)
	pid.feedbackTime = t
	pid.calcError(processVariable)
	pid.calcP()
	pid.calcI()
	pid.calcD()
	return utils.Max(pid.pControlVariable+pid.iControlVariable+pid.dControlVariable, pid.maxWeightedSum)
}

func (pid *PIDControl) calcError(processVariable float64) {
	pid.prevErrorValue = pid.errorValue
	pid.errorValue = utils.Max(pid.setPoint-processVariable, pid.maxError)
}

func (pid *PIDControl) SetSetPoint(setPoint float64) {
	pid.setPoint = setPoint
}

func (pid *PIDControl) calcP() {
	pid.pControlVariable = pid.pGain * pid.errorValue
}

func (pid *PIDControl) calcI() {
	pid.iControlVariable += pid.iGain * pid.errorValue * pid.dt.Seconds()
	pid.iControlVariable = utils.Max(pid.iControlVariable, pid.maxIntegrationValue)
}

func (pid *PIDControl) calcD() {
	if pid.dt < MIN_DT {
		pid.dControlVariable = 0
	} else {
		pid.dControlVariable = pid.dGain * (pid.errorValue - pid.prevErrorValue) / pid.dt.Seconds()
	}
}

func (pid *PIDControl) Reset() {
	pid.iControlVariable = 0
}

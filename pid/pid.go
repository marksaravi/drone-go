package pid

import (
	"fmt"
	"math"
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
	CalibrationMode     bool
	CalibrationIncP     float64
	CalibrationIncI     float64
	CalibrationIncD     float64
}

type average interface {
	AddValue(v float64) float64
	Average() float64
}

type PIDControl struct {
	pGain               float64
	iGain               float64
	dGain               float64
	outDataPerInputData int
	pControlVariable    average
	dControlVariable    average
	maxError            float64
	maxIntegrationValue float64
	minProcessVariable  float64
	maxWeightedSum      float64
	setPoint            float64
	feedbackTime        time.Time
	dt                  time.Duration
	errorValue          float64
	prevErrorValue      float64
	iControlVariable    float64
	weightedSum         float64
}

func NewPIDControl(settings PIDSettings, outDataPerInputData int) *PIDControl {
	pid := &PIDControl{
		pGain:               settings.PGain,
		iGain:               settings.IGain,
		dGain:               settings.DGain,
		outDataPerInputData: outDataPerInputData,
		pControlVariable:    utils.NewAverage[float64](outDataPerInputData),
		dControlVariable:    utils.NewAverage[float64](outDataPerInputData),
		maxError:            settings.MaxError,
		maxIntegrationValue: settings.MaxIntegrationValue,
		minProcessVariable:  settings.MinProcessVariable,
		maxWeightedSum:      settings.MaxWeightedSum,
		setPoint:            0,
		feedbackTime:        time.Now().Add(time.Second * 32000000),
		dt:                  0,
		errorValue:          0,
		prevErrorValue:      0,
		iControlVariable:    0,
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
	return max(pid.pControlVariable.Average()+pid.iControlVariable+pid.dControlVariable.Average(), pid.maxWeightedSum)
}

func (pid *PIDControl) calcError(processVariable float64) {
	pid.prevErrorValue = pid.errorValue
	pid.errorValue = utils.Max(pid.setPoint-processVariable, pid.maxError)
}

func (pid *PIDControl) SetSetPoint(setPoint float64) {
	pid.setPoint = setPoint
}

func (pid *PIDControl) calcP() {
	pControlVariable := pid.pGain * pid.errorValue
	pid.pControlVariable.AddValue(pControlVariable)
}

func (pid *PIDControl) calcI() {
	pid.iControlVariable += pid.iGain * pid.errorValue * pid.dt.Seconds()
	pid.iControlVariable = max(pid.iControlVariable, pid.maxIntegrationValue)
}

func (pid *PIDControl) calcD() {
	dControlVariable := float64(0)
	if pid.dt >= MIN_DT {
		dControlVariable = pid.dGain * (pid.errorValue - pid.prevErrorValue) / pid.dt.Seconds()
	}
	pid.dControlVariable.AddValue(dControlVariable)
}

func (pid *PIDControl) Reset() {
	pid.iControlVariable = 0
}

func (pid *PIDControl) UpdateGainP(v float64) {
	pid.pGain += v
}

func (pid *PIDControl) UpdateGainI(v float64) {
	pid.iGain += v
}

func (pid *PIDControl) UpdateGainD(v float64) {
	pid.dGain += v
}

func (pid *PIDControl) GainP() float64 {
	return pid.pGain
}

func (pid *PIDControl) GainI() float64 {
	return pid.iGain
}

func (pid *PIDControl) GainD() float64 {
	return pid.dGain
}

func max(v, maxValue float64) float64 {
	if math.Abs(v) < maxValue {
		return v
	}
	fmt.Println(v)
	if v < 0 {
		return -maxValue
	}
	return maxValue
}

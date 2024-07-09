package pid

import (
	"time"
)

type PIDSettings struct {
	PGain   float64
	IGain   float64
	DGain   float64
	Enabled bool
}

type PIDControl struct {
	id               string
	pGain            float64
	iGain            float64
	dGain            float64
	pControlVariable float64
	iControlVariable float64
	dControlVariable float64
	setPoint         float64
	prevFeedbackTime time.Time
	dt               time.Duration
	errorValue       float64
	prevErrorValue   float64
	enabled          bool
}

func NewPIDControl(id string, settings PIDSettings, outDataPerInputData int) *PIDControl {
	pid := &PIDControl{
		id:               id,
		pGain:            settings.PGain,
		iGain:            settings.IGain,
		dGain:            settings.DGain,
		pControlVariable: 0,
		iControlVariable: 0,
		dControlVariable: 0,
		setPoint:         0,
		prevFeedbackTime: time.Now().Add(time.Second * 32000000),
		errorValue:       0,
		prevErrorValue:   0,
		enabled:          settings.Enabled,
	}
	return pid
}

func (pid *PIDControl) CalculateControlVariable(processVariable float64, t time.Time) float64 {
	if !pid.enabled {
		return 0
	}
	pid.dt = t.Sub(pid.prevFeedbackTime)
	pid.prevFeedbackTime = t
	pid.calcError(processVariable)
	pid.calcP()
	pid.calcI()
	pid.calcD()
	return pid.pControlVariable + pid.iControlVariable + pid.dControlVariable
}

func (pid *PIDControl) calcError(processVariable float64) {
	pid.prevErrorValue = pid.errorValue
	pid.errorValue = pid.setPoint - processVariable
}

func (pid *PIDControl) SetSetPoint(setPoint float64) {
	pid.setPoint = setPoint
}

func (pid *PIDControl) calcP() {
	pid.pControlVariable = pid.pGain * pid.errorValue
}

func (pid *PIDControl) calcI() {
	pid.iControlVariable += pid.iGain * pid.errorValue * pid.dt.Seconds()
}

func (pid *PIDControl) calcD() {
	pid.dControlVariable = pid.dGain * (pid.errorValue - pid.prevErrorValue) / pid.dt.Seconds()
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

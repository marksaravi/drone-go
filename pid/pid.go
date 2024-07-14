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
	id             string
	kP             float64
	kI             float64
	kD             float64
	integralValue  float64
	setPoint       float64
	prevTime       time.Time
	prevErrorValue float64
	enabled        bool
}

func NewPIDControl(id string, settings PIDSettings, outDataPerInputData int) *PIDControl {
	pid := &PIDControl{
		id:             id,
		kP:             settings.PGain,
		kI:             settings.IGain,
		kD:             settings.DGain,
		setPoint:       0,
		prevTime:       time.Now().Add(time.Second * 32000000),
		prevErrorValue: 0,
		enabled:        settings.Enabled,
	}
	return pid
}

func (pid *PIDControl) CalcProcessValue(measuredValue float64, t time.Time, processOffset float64, sign int) (float64, float64) {
	u := pid.calcProcessValue(measuredValue, t)
	return processOffset + float64(sign)*u, processOffset + float64(-sign)*u
}

func (pid *PIDControl) calcProcessValue(measuredValue float64, t time.Time) float64 {
	if !pid.enabled {
		return 0
	}
	errorValue := measuredValue - pid.setPoint
	dErrorValue := errorValue - pid.prevErrorValue
	pid.prevErrorValue = errorValue
	dt := t.Sub(pid.prevTime)
	pid.prevTime = t

	p := pid.kP * errorValue
	pid.integralValue = errorValue * dt.Seconds()
	i := pid.kI * pid.integralValue
	deDt := dErrorValue / dt.Seconds()
	d := pid.kD * deDt
	u := p + i + d
	return u
}

func (pid *PIDControl) ResetI() {
	pid.integralValue = 0
}

func (pid *PIDControl) SetSetPoint(setPoint float64) {
	pid.setPoint = setPoint
}

func (pid *PIDControl) UpdateGainP(v float64) {
	pid.kP += v
}

func (pid *PIDControl) UpdateGainI(v float64) {
	pid.kI += v
}

func (pid *PIDControl) UpdateGainD(v float64) {
	pid.kD += v
}

func (pid *PIDControl) GainP() float64 {
	return pid.kP
}

func (pid *PIDControl) GainI() float64 {
	return pid.kI
}

func (pid *PIDControl) GainD() float64 {
	return pid.kD
}

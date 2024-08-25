package pid

import (
	"fmt"
	"time"
)

type PIDConfigs struct {
	Id                  string  `json:"id"`
	PGain               float64 `json:"p-gain"`
	IGain               float64 `json:"i-gain"`
	DGain               float64 `json:"d-gain"`
	MaxRotationError    float64 `json:"max-rot-error"`
	MaxIntegrationValue float64 `json:"max-i-value"`
	MaxWeightedSum      float64 `json:"max-weighted-sum"`
	CalibrationMode     bool    `json:"calibration-mode"`
	CalibrationIncP     float64 `json:"calibration-p-inc"`
	CalibrationIncI     float64 `json:"calibration-i-inc"`
	CalibrationIncD     float64 `json:"calibration-d-inc"`
}

type PIDControl struct {
	id             string
	settings       PIDConfigs
	integralValue  float64
	setPoint       float64
	prevTime       time.Time
	prevErrorValue float64
}

func NewPIDControl(id string, settings PIDConfigs) *PIDControl {
	pid := &PIDControl{
		id:             settings.Id,
		settings:       settings,
		setPoint:       0,
		prevTime:       time.Now().Add(time.Second * 32000000),
		prevErrorValue: 0,
	}
	return pid
}

func (pid *PIDControl) CalcOutput(measuredValue float64, t time.Time, processOffset float64, sign int) (float64, float64) {
	u := pid.calcProcessValue(measuredValue, t)
	return processOffset + float64(sign)*u, processOffset - float64(sign)*u
}

func (pid *PIDControl) calcProcessValue(measuredValue float64, t time.Time) float64 {
	errorValue := measuredValue - pid.setPoint
	dErrorValue := errorValue - pid.prevErrorValue
	pid.prevErrorValue = errorValue
	dt := t.Sub(pid.prevTime)
	pid.prevTime = t

	p := pid.settings.PGain * errorValue

	pid.integralValue = errorValue * dt.Seconds()
	i := pid.settings.IGain * pid.integralValue

	deDt := dErrorValue / dt.Seconds()
	d := pid.settings.DGain * deDt

	u := p + i + d
	return u
}

func (pid *PIDControl) IsCalibrationEnabled() bool {
	return pid.settings.CalibrationMode
}

func (pid *PIDControl) Calibrate(t rune, inc bool) {
	if !pid.settings.CalibrationMode {
		return
	}
	switch t {
	case 'p':
		pid.settings.PGain = updateGain(pid.settings.PGain, pid.settings.CalibrationIncP, inc)
	case 'i':
		pid.settings.IGain = updateGain(pid.settings.IGain, pid.settings.CalibrationIncI, inc)
	case 'd':
		pid.settings.DGain = updateGain(pid.settings.DGain, pid.settings.CalibrationIncD, inc)
	}
	fmt.Printf("%s -> %c:%8.2f, %8.2f, %8.2f\n ", pid.id, t, pid.settings.PGain, pid.settings.IGain, pid.settings.DGain)
}

func updateGain(v, c float64, inc bool) float64 {
	if inc {
		return v + c
	}
	if v-c > 0 {
		return v - c
	}
	return 0
}

func (pid *PIDControl) ResetI() {
	pid.integralValue = 0
}

func (pid *PIDControl) SetTargetRotation(setPoint float64) {
	pid.setPoint = setPoint
}

func (pid *PIDControl) UpdateGainP(v float64) {
	pid.settings.PGain += v
}

func (pid *PIDControl) UpdateGainI(v float64) {
	pid.settings.IGain += v
}

func (pid *PIDControl) UpdateGainD(v float64) {
	pid.settings.DGain += v
}

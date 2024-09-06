package pid

import (
	"fmt"
	"time"

	"github.com/marksaravi/drone-go/utils"
)

type PIDConfigs struct {
	Id                  string  `json:"id"`
	PGain               float64 `json:"p-gain"`
	IGain               float64 `json:"i-gain"`
	DGain               float64 `json:"d-gain"`
	Direction           float64 `json:"direction"`
	MaxRotationError    float64 `json:"max-rot-error"`
	MaxIntegrationValue float64 `json:"max-i-value"`
	MaxDiffValue        float64 `json:"max-d-value"`
	CalibrationMode     bool    `json:"calibration-mode"`
	CalibrationIncP     float64 `json:"calibration-p-inc"`
	CalibrationIncI     float64 `json:"calibration-i-inc"`
	CalibrationIncD     float64 `json:"calibration-d-inc"`
}

type PIDControl struct {
	id            string
	settings      PIDConfigs
	integralValue float64
	setPoint      float64
	prevTime      time.Time
	prevGyroRot   float64
}

func NewPIDControl(id string, settings PIDConfigs) *PIDControl {
	pid := &PIDControl{
		id:          settings.Id,
		settings:    settings,
		setPoint:    0,
		prevTime:    time.Now().Add(time.Second * 32000000),
		prevGyroRot: 0,
	}
	return pid
}

func (pid *PIDControl) CalcOutput(rot, gyroRot float64, t time.Time) float64 {
	return pid.calcProcessValue(rot, gyroRot, t) * float64(pid.settings.Direction)
}

func (pid *PIDControl) calcProcessValue(rot, gyroRot float64, t time.Time) float64 {
	errorValue := utils.SignedMax(rot-pid.setPoint, pid.settings.MaxRotationError)
	dRot := gyroRot - pid.prevGyroRot
	pid.prevGyroRot = gyroRot
	dt := t.Sub(pid.prevTime)
	pid.prevTime = t

	p := pid.settings.PGain * errorValue

	pid.integralValue = utils.SignedMax(pid.integralValue+errorValue*dt.Seconds()*pid.settings.IGain, pid.settings.MaxIntegrationValue)

	d := pid.settings.DGain * dRot / dt.Seconds()

	u := p + pid.integralValue + d
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
	fmt.Printf("%s -> %c:%8.4f, %8.4f, %8.4f\n ", pid.id, t, pid.settings.PGain, pid.settings.IGain, pid.settings.DGain)
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

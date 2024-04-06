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
}

type PIDControl struct {
	pGain               float64
	iGain               float64
	dGain               float64
	maxError            float64
	maxIntegrationValue float64
	minProcessVariable  float64
	maxWeightedSum      float64
	maxOutput           float64
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

func (pid *PIDControl) CalculateControlVariable(processVariable float64, t time.Time) {
	if t.Sub(pid.feedbackTime) < 0 {
		pid.feedbackTime = t
	}
	pid.dt = t.Sub(pid.feedbackTime)
	pid.feedbackTime = t
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
	if math.Abs(pid.iControlVariable) > pid.maxIntegrationValue {
		if pid.iControlVariable > 0 {
			pid.iControlVariable = pid.maxIntegrationValue
		} else {
			pid.iControlVariable = -pid.maxIntegrationValue
		}
	}
}

func (pid *PIDControl) calcD() {
	if pid.dt < MIN_DT {
		pid.dControlVariable = 0
	} else {
		pid.dControlVariable = pid.dGain * (pid.errorValue - pid.prevErrorValue) / pid.dt.Seconds()
	}
}

func (pid *PIDControl) applyMaxRotationError(value, prevValue float64) float64 {
	v := value - prevValue
	if v < -pid.maxRotError {
		return -pid.maxRotError
	}
	if v > pid.maxRotError {
		return pid.maxRotError
	}
	return v
}

func (pid *PIDControl) resetIValues() {
	pid.arm_0_2_i_value = 0
	pid.arm_1_3_i_value = 0
}

func (pid *PIDControl) memoFeedback() {
	pid.prevFeedback = pid.feedback
	pid.prevErrorValue = pid.errorValue
}
func (pid *PIDControl) calcThrottlesFlightMode() {
	pid.applyP()
	pid.applyI()
	pid.applyD()
	for i := 0; i < 4; i++ {
		pid.throttles[i] = pid.pThrottles[i] + pid.iThrottles[i] + pid.dThrottles[i] + pid.throttle
	}
}

func (pid *PIDControl) calcThrottlesLowPowerMode() {
	pid.resetIValues()
	for i := 0; i < 4; i++ {
		pid.throttles[i] = pid.throttle
	}
}

func (pid *PIDControl) CalcESCThrottles() {
	pid.calcErrorValues()
	if pid.throttle >= pid.minFlightThrottle {
		pid.calcThrottlesFlightMode()
	} else {
		pid.calcThrottlesLowPowerMode()
	}
	pid.memoFeedback()
}

var throttleDisplay = utils.WithDataPerSecond(5)

func (pid *PIDControl) GetThrottles() []float64 {
	if throttleDisplay.IsTime() {
		fmt.Printf("%6.1f %6.1f %6.1f %6.1f %6.1f %6.1f %6.1f %6.1f %6.1f %6.1f %6.1f\n", pid.throttle, pid.throttles[0], pid.throttles[1], pid.throttles[2], pid.throttles[3], pid.arm_0_2_d_rotError, pid.arm_1_3_d_rotError, pid.iThrottles[0], pid.iThrottles[1], pid.iThrottles[2], pid.iThrottles[3])
	}
	return pid.throttles
}

// if rotDisplay.IsTime() {
// 	fmt.Printf("%6.1f,%6.1f,%6.1f,%6.1f,%6.1f,%6.1f,%6.1f,%6.1f,%6.1f,%6.1f,%6.1f,%6.1f\n",
// 		pid.feedback.Roll, pid.feedback.Pitch,
// 		pid.setPoint.Roll, pid.setPoint.Pitch,
// 		pid.errorValue.Roll, pid.errorValue.Pitch,
// 		pid.arm_0_2_rotError, pid.arm_1_3_rotError,
// 		pid.throttles[0], pid.throttles[1], pid.throttles[2], pid.throttles[3])
// }

package drone

import (
	"fmt"
	"time"

	"github.com/marksaravi/drone-go/devices/imu"
	"github.com/marksaravi/drone-go/utils"
)

type PIDControl struct {
	pGain               float64
	iGain               float64
	dGain               float64
	maxRotError         float64
	maxIntegrationValue float64

	rotations          imu.Rotations
	prevRotations      imu.Rotations
	targetRotations    imu.Rotations
	rotationsError     imu.Rotations
	prevRotationsError imu.Rotations
	arm_0_2_rotError   float64
	arm_1_3_rotError   float64
	arm_0_2_i_value    float64
	arm_1_3_i_value    float64

	throttle     float64
	prevThrottle float64
	pThrottles   []float64
	iThrottles   []float64
	dThrottles   []float64
	throttles    []float64
}

func NewPIDControl(pidCongigs PIDConfigs) *PIDControl {
	fmt.Println("PID: ", pidCongigs)
	return &PIDControl{
		pGain:               pidCongigs.P,
		iGain:               pidCongigs.I,
		dGain:               pidCongigs.D,
		maxRotError:         pidCongigs.MaxRotationError,
		maxIntegrationValue: pidCongigs.MaxIntegrationValue,
		throttle:            0,
		prevThrottle:        0,
		pThrottles:          make([]float64, 4),
		iThrottles:          make([]float64, 4),
		dThrottles:          make([]float64, 4),
		throttles:           make([]float64, 4),
		rotations: imu.Rotations{
			Roll:  0,
			Pitch: 0,
			Yaw:   0,
			Time:  time.Now(),
		},
		targetRotations: imu.Rotations{
			Roll:  0,
			Pitch: 0,
			Yaw:   0,
		},
		prevRotations: imu.Rotations{
			Roll:  0,
			Pitch: 0,
			Yaw:   0,
			Time:  time.Now(),
		},
		rotationsError: imu.Rotations{
			Roll:  0,
			Pitch: 0,
			Yaw:   0,
		},
		prevRotationsError: imu.Rotations{
			Roll:  0,
			Pitch: 0,
			Yaw:   0,
		},
	}
}

var rotDisplay = utils.WithDataPerSecond(5)

func (pid *PIDControl) Reset(prevThrottle float64) {
	pid.arm_0_2_i_value = 0
	pid.arm_1_3_i_value = 0
	pid.prevThrottle = prevThrottle
	fmt.Println("FLIGHT THROTTLE STATE", pid.prevThrottle)

}

func (pid *PIDControl) SetRotations(rotattions imu.Rotations) {
	pid.rotations = rotattions
}

func (pid *PIDControl) SetTargetRotations(rotattions imu.Rotations) {
	pid.targetRotations = rotattions
}

func (pid *PIDControl) SetThrottle(throttle float64) {
	const K = float64(0.45)
	pid.throttle = throttle*K + pid.prevThrottle*(1-K)
	pid.prevThrottle = pid.throttle
}

func (pid *PIDControl) applyP() {
	gain_0_2 := pid.pGain * pid.arm_0_2_rotError
	gain_1_3 := pid.pGain * pid.arm_1_3_rotError
	pid.pThrottles[0] = gain_0_2
	pid.pThrottles[1] = -gain_1_3
	pid.pThrottles[2] = -gain_0_2
	pid.pThrottles[3] = gain_1_3
}

func (pid *PIDControl) applyI() {
}

func (pid *PIDControl) applyD() {
}

func (pid *PIDControl) calcRotationsErrors() {
	pid.rotationsError.Roll = pid.calcRotationsError(pid.targetRotations.Roll, pid.rotations.Roll)
	pid.rotationsError.Pitch = pid.calcRotationsError(pid.targetRotations.Pitch, pid.rotations.Pitch)
	pid.arm_0_2_rotError = (pid.rotationsError.Pitch + pid.rotationsError.Roll) / 2
	pid.arm_1_3_rotError = (pid.rotationsError.Pitch - pid.rotationsError.Roll) / 2
	if rotDisplay.IsTime() {
		fmt.Printf("%6.1f,%6.1f,%6.1f,%6.1f,%6.1f,%6.1f,%6.1f,%6.1f\n",
			pid.rotations.Roll, pid.rotations.Pitch,
			pid.targetRotations.Roll, pid.targetRotations.Pitch,
			pid.rotationsError.Roll, pid.rotationsError.Pitch,
			pid.arm_0_2_rotError, pid.arm_1_3_rotError)
	}
}

func (pid *PIDControl) calcRotationsError(value, prevValue float64) float64 {
	v := value - prevValue
	if v < -pid.maxRotError {
		return -pid.maxRotError
	}
	if v > pid.maxRotError {
		return pid.maxRotError
	}
	return v
}

func (pid *PIDControl) CalcThrottles() []float64 {
	pid.calcRotationsErrors()
	pid.applyP()
	pid.applyI()
	pid.applyD()
	for i := 0; i < 4; i++ {
		pid.throttles[i] = pid.pThrottles[i] + pid.iThrottles[i] + pid.dThrottles[i] + pid.throttle
	}
	pid.prevRotations = pid.rotations
	pid.prevRotationsError = pid.rotationsError
	return pid.throttles
}

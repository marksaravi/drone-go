package drone

import (
	"fmt"
	"time"

	"github.com/marksaravi/drone-go/devices/imu"
	"github.com/marksaravi/drone-go/utils"
)

type PIDControl struct {
	pGain float64
	iGain float64
	dGain float64

	rotations       imu.Rotations
	prevRotations   imu.Rotations
	targetRotations imu.Rotations
	throttle        float64
	prevThrottle    float64

	pThrottles []float64
	iThrottles []float64
	dThrottles []float64
	throttles  []float64
}

func NewPIDControl(pGain, iGain, dGain float64) *PIDControl {
	return &PIDControl{
		pGain:        pGain,
		iGain:        iGain,
		dGain:        dGain,
		throttle:     0,
		prevThrottle: 0,
		pThrottles:   make([]float64, 4),
		iThrottles:   make([]float64, 4),
		dThrottles:   make([]float64, 4),
		throttles:    make([]float64, 4),
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
	}
}

var rotDisplay = utils.WithDataPerSecond(5)

func (pid *PIDControl) SetRotations(rotattions imu.Rotations) {
	pid.rotations = rotattions
	if rotDisplay.IsTime() {
		fmt.Printf("%6.1f,%6.1f,%6.1f,%v\n", pid.rotations.Roll, pid.rotations.Pitch, pid.rotations.Yaw, pid.rotations.Time)
	}
}

func (pid *PIDControl) SetTargetRotations(rotattions imu.Rotations) {
	pid.targetRotations = rotattions
}

func (pid *PIDControl) CalcThrottles() []float64 {
	pid.calcRotationsDiff()
	pid.applyP()
	pid.applyI()
	pid.applyD()
	for i := 0; i < 4; i++ {
		pid.throttles[i] = pid.pThrottles[i] + pid.iThrottles[i] + pid.dThrottles[i] + pid.throttle
	}
	pid.prevRotations = pid.rotations
	return pid.throttles
}

func (pid *PIDControl) SetThrottle(throttle float64) {
	const K = float64(0.45)
	pid.throttle = throttle*K + pid.prevThrottle*(1-K)
	pid.prevThrottle = pid.throttle
}

func (pid *PIDControl) applyP() []float64 {
	return pid.pThrottles
}

func (pid *PIDControl) applyI() []float64 {
	return pid.iThrottles
}

func (pid *PIDControl) applyD() []float64 {
	return pid.dThrottles
}

func (pid *PIDControl) calcRotationsDiff() {

}

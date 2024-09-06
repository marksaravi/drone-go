package drone

import (
	"github.com/marksaravi/drone-go/devices/imu"
	"github.com/marksaravi/drone-go/pid"
	"github.com/marksaravi/drone-go/utils"
)

const (
	MOTORS_OFF         = false
	MOTORS_ONN         = true
	THROTTLE_HYSTERSYS = 0.5
)

type FlightControl struct {
	arm_0_2_pid           *pid.PIDControl
	arm_1_3_pid           *pid.PIDControl
	yaw_pid               *pid.PIDControl
	throttle              float64
	outputThrottles       []float64
	maxThrottle           float64
	maxOutputThrottle     float64
	escs                  escs
	throttleLowPassFilter float64
}

func NewFlightControl(
	escs escs,
	maxThrottle float64,
	maxOutputThrottle float64,
	arm_0_2_pid *pid.PIDControl,
	arm_1_3_pid *pid.PIDControl,
	yaw_pid *pid.PIDControl,
) *FlightControl {
	fc := &FlightControl{
		arm_0_2_pid:           arm_0_2_pid,
		arm_1_3_pid:           arm_1_3_pid,
		yaw_pid:               yaw_pid,
		throttleLowPassFilter: 0.45,
		throttle:              0,
		maxThrottle:           maxThrottle,
		maxOutputThrottle:     maxOutputThrottle,
		outputThrottles:       make([]float64, 4),
		escs:                  escs,
	}
	return fc
}

func transformRollPitch(roll, pitch float64) (float64, float64) {
	return (pitch + roll) / 2, (pitch - roll) / 2
}

func (fc *FlightControl) calcOutputThrottles(rotattions imu.Rotations, gyroRotattions imu.Rotations) {
	arm_0_2_rotation, arm_1_3_rotation := transformRollPitch(rotattions.Roll, rotattions.Pitch)
	arm_0_2_grotation, arm_1_3_grotation := transformRollPitch(gyroRotattions.Roll, gyroRotattions.Pitch)

	arm_o_2_pid := fc.arm_0_2_pid.CalcOutput(arm_0_2_rotation, arm_0_2_grotation, rotattions.Time)
	motor_0_output_throttle := fc.throttle + arm_o_2_pid
	motor_2_output_throttle := fc.throttle - arm_o_2_pid

	arm_1_3_pid := fc.arm_1_3_pid.CalcOutput(arm_1_3_rotation, arm_1_3_grotation, rotattions.Time)
	motor_1_output_throttle := fc.throttle + arm_1_3_pid
	motor_3_output_throttle := fc.throttle - arm_1_3_pid

	fc.outputThrottles[0] = motor_0_output_throttle
	fc.outputThrottles[2] = motor_2_output_throttle

	fc.outputThrottles[1] = motor_1_output_throttle * 0
	fc.outputThrottles[3] = motor_3_output_throttle * 0
}

func (fc *FlightControl) setTargetRotations(rotattions imu.Rotations) {
	arm_0_2_rotation, arm_1_3_rotation := transformRollPitch(rotattions.Roll, rotattions.Pitch)

	fc.arm_0_2_pid.SetTargetRotation(arm_0_2_rotation)
	fc.arm_1_3_pid.SetTargetRotation(arm_1_3_rotation)
}

func (fc *FlightControl) setThrottle(throttle float64) {
	fc.throttle = throttle
}

func (fc *FlightControl) getThrottle() float64 {
	return fc.throttle
}

func (fc *FlightControl) applyThrottles() {
	for i := 0; i < len(fc.outputThrottles); i++ {
		fc.outputThrottles[i] = utils.SignedMax(fc.outputThrottles[i], fc.maxOutputThrottle)
	}
	fc.escs.SetThrottles(fc.outputThrottles)
}

func (fc *FlightControl) turnOnMotors(motorsOn bool) {
	if motorsOn {
		fc.escs.On()
	} else {
		fc.escs.Off()
	}
}

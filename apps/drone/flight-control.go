package drone

import (
	"github.com/marksaravi/drone-go/devices/imu"
	"github.com/marksaravi/drone-go/pid"
)

const (
	MOTORS_OFF         = false
	MOTORS_ONN         = true
	THROTTLE_HYSTERSYS = 0.5
)

type FlightControl struct {
	arm_0_2_PID     *pid.PIDControl
	arm_1_3_PID     *pid.PIDControl
	yawPID          *pid.PIDControl
	calibrationMode bool
	calibrationIncP float64
	calibrationIncI float64
	calibrationIncD float64

	throttle        float64
	outputThrottles []float64
	maxThrottle     float64

	escs                  escs
	throttleLowPassFilter float64
}

func NewFlightControl(escs escs, maxThrottle float64, pidSettings pid.PIDSettings, escsDataPerImuData int) *FlightControl {
	fc := &FlightControl{
		arm_0_2_PID:           pid.NewPIDControl("0_2", pidSettings, escsDataPerImuData),
		arm_1_3_PID:           pid.NewPIDControl("1_3", pidSettings, escsDataPerImuData),
		throttleLowPassFilter: 0.45,
		throttle:              0,
		maxThrottle:           maxThrottle,

		outputThrottles: make([]float64, 4),
		escs:            escs,
	}
	return fc
}

func transformRollPitch(roll, pitch float64) (float64, float64) {
	return (pitch + roll) / 2, (pitch - roll) / 2
}

func (fc *FlightControl) calcOutputThrottles(rotattions imu.Rotations) {
	arm_0_2_rotation, _ := transformRollPitch(rotattions.Roll, rotattions.Pitch)

	motor_0_output_throttle, motor_2_output_throttle := fc.arm_0_2_PID.CalcOutput(arm_0_2_rotation, rotattions.Time, fc.throttle, 1)
	// motor_1_output_throttle, motor_3_output_throttle := fc.arm_1_3_PID.CalcOutput(arm_1_3_rotation, rotattions.Time, fc.throttle, 1)

	fc.outputThrottles[0] = motor_0_output_throttle
	fc.outputThrottles[2] = motor_2_output_throttle

	fc.outputThrottles[1] = 0 //motor_1_output_throttle
	fc.outputThrottles[3] = 0 //motor_3_output_throttle
}

func (fc *FlightControl) setTargetRotations(rotattions imu.Rotations) {
	arm_0_2_rotation, arm_1_3_rotation := transformRollPitch(rotattions.Roll, rotattions.Pitch)

	fc.arm_0_2_PID.SetTargetRotation(arm_0_2_rotation)
	fc.arm_1_3_PID.SetTargetRotation(arm_1_3_rotation)
}

func (fc *FlightControl) setThrottle(throttle float64) {
	fc.throttle = throttle
}

func (fc *FlightControl) getThrottle() float64 {
	return fc.throttle
}

func (fc *FlightControl) applyThrottles() {
	fc.escs.SetThrottles(fc.outputThrottles)
}

func (fc *FlightControl) turnOnMotors(motorsOn bool) {
	if motorsOn {
		fc.escs.On()
	} else {
		fc.escs.Off()
	}
}

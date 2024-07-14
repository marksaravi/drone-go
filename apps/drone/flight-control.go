package drone

import (
	"time"

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

	throttle     float64
	pidThrottles []float64
	maxThrottle  float64

	escs                  escs
	motorsArmingTime      time.Time
	throttleLowPassFilter float64
}

func NewFlightControl(escs escs, maxThrottle float64, pidSettings pid.PIDSettings, escsDataPerImuData int) *FlightControl {
	fc := &FlightControl{
		arm_0_2_PID:           pid.NewPIDControl("0_2", pidSettings, escsDataPerImuData),
		arm_1_3_PID:           pid.NewPIDControl("1_3", pidSettings, escsDataPerImuData),
		throttleLowPassFilter: 0.45,
		throttle:              0,
		maxThrottle:           maxThrottle,

		pidThrottles: make([]float64, 4),
		escs:         escs,
	}
	return fc
}

func transformRollPitch(roll, pitch float64) (float64, float64) {
	return (pitch + roll) / 2, (pitch - roll) / 2
}

func (fc *FlightControl) calcThrottles(rotattions imu.Rotations) {
	arm_0_2_rotation, arm_1_3_rotation := transformRollPitch(rotattions.Roll, rotattions.Pitch)
	motor_0_Throttle, motor_2_Throttle := fc.arm_0_2_PID.CalcProcessValue(arm_0_2_rotation, rotattions.Time, fc.throttle, 1)
	motor_1_Throttle, motor_3_Throttle := fc.arm_1_3_PID.CalcProcessValue(arm_1_3_rotation, rotattions.Time, fc.throttle, 1)

	fc.pidThrottles[0] = motor_0_Throttle
	fc.pidThrottles[2] = motor_2_Throttle

	fc.pidThrottles[1] = motor_1_Throttle
	fc.pidThrottles[3] = motor_3_Throttle
}

func (fc *FlightControl) setTargetRotations(rotattions imu.Rotations) {
	arm_0_2_rotation, arm_1_3_rotation := transformRollPitch(rotattions.Roll, rotattions.Pitch)

	fc.arm_0_2_PID.SetSetPoint(arm_0_2_rotation)
	fc.arm_1_3_PID.SetSetPoint(arm_1_3_rotation)
}

func (fc *FlightControl) setThrottle(throttle float64) {
	fc.throttle = throttle
}

func (fc *FlightControl) applyThrottles() {
	fc.escs.SetThrottles(fc.pidThrottles)
}

func (fc *FlightControl) turnOnMotors(motorsOn bool) {
	if motorsOn {
		fc.escs.On()
	} else {
		fc.escs.Off()
	}
}

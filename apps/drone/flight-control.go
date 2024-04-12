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
	arm_0_2_PID *pid.PIDControl
	arm_1_3_PID *pid.PIDControl

	throttle              float64
	escs                  escs
	motorsArmingTime      time.Time
	throttleLowPassFilter float64
}

func NewFlightControl(escs escs, minFlightThrottle, maxThrottle float64, pidSettings pid.PIDSettings) *FlightControl {
	fc := &FlightControl{
		arm_0_2_PID:           pid.NewPIDControl(pidSettings),
		arm_1_3_PID:           pid.NewPIDControl(pidSettings),
		throttleLowPassFilter: 0.45,
		throttle:              0,
		escs:                  escs,
	}
	fc.arm_0_2_PID.Initiate()
	fc.arm_1_3_PID.Initiate()

	fc.turnOnMotors(false)
	return fc
}

func transformRollPitch(roll, pitch float64) (float64, float64) {
	return (pitch + roll) / 2, (pitch - roll) / 2
}

func (fc *FlightControl) SetRotations(rotattions imu.Rotations) {
	arm_0_2_rotation, arm_1_3_rotation := transformRollPitch(rotattions.Roll, rotattions.Pitch)

	fc.arm_0_2_PID.CalculateControlVariable(arm_0_2_rotation, rotattions.Time)
	fc.arm_1_3_PID.CalculateControlVariable(arm_1_3_rotation, rotattions.Time)
}

func (fc *FlightControl) SetTargetRotations(rotattions imu.Rotations) {
	arm_0_2_rotation, arm_1_3_rotation := transformRollPitch(rotattions.Roll, rotattions.Pitch)

	fc.arm_0_2_PID.SetSetPoint(arm_0_2_rotation)
	fc.arm_1_3_PID.SetSetPoint(arm_1_3_rotation)
}

func (fc *FlightControl) SetThrottle(throttle float64) {
	fc.throttle = fc.throttle*(1-fc.throttleLowPassFilter) + fc.throttleLowPassFilter*throttle
}

func (fc *FlightControl) ApplyThrottlesToESCs() {
	if time.Since(fc.motorsArmingTime) >= 0 {
		fc.escs.SetThrottles([]float64{fc.throttle, fc.throttle, fc.throttle, fc.throttle})
	} else {
		fc.escs.SetThrottles([]float64{0, 0, 0, 0})
	}
}

func (fc *FlightControl) turnOnMotors(motorsOn bool) {
	if motorsOn && fc.throttle < 2 {
		fc.setArmingTime(true)
		fc.escs.On()
	} else if !motorsOn {
		fc.setArmingTime(false)
		fc.escs.Off()
	}
}

func (fc *FlightControl) setArmingTime(on bool) {
	if on {
		fc.motorsArmingTime = time.Now().Add(time.Second * 3)
	} else {
		fc.motorsArmingTime = time.Now().Add(time.Second * 32000000)
	}
}

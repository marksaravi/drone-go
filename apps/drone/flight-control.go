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
		calibrationMode:       pidSettings.CalibrationMode,
		calibrationIncP:       pidSettings.CalibrationIncP,
		calibrationIncI:       pidSettings.CalibrationIncI,
		calibrationIncD:       pidSettings.CalibrationIncD,
		throttleLowPassFilter: 0.45,
		throttle:              0,
		maxThrottle:           maxThrottle,

		pidThrottles: make([]float64, 4),
		escs:         escs,
	}
	fc.resetPIDs()
	return fc
}

func transformRollPitch(roll, pitch float64) (float64, float64) {
	return (pitch + roll) / 2, (pitch - roll) / 2
}

func (fc *FlightControl) resetPIDs() {
	fc.arm_0_2_PID.Reset()
	fc.arm_1_3_PID.Reset()

}

func (fc *FlightControl) SetRotations(rotattions imu.Rotations) {
	arm_0_2_rotation, arm_1_3_rotation := transformRollPitch(rotattions.Roll, rotattions.Pitch)
	arm_0_2_controlVariable := fc.arm_0_2_PID.CalculateControlVariable(arm_0_2_rotation, rotattions.Time)
	arm_1_3_controlVariable := fc.arm_1_3_PID.CalculateControlVariable(arm_1_3_rotation, rotattions.Time)

	fc.pidThrottles[0] = fc.throttle + arm_0_2_controlVariable
	fc.pidThrottles[2] = fc.throttle - arm_0_2_controlVariable

	fc.pidThrottles[1] = fc.throttle - arm_1_3_controlVariable
	fc.pidThrottles[3] = fc.throttle + arm_1_3_controlVariable
}

func (fc *FlightControl) SetTargetRotations(rotattions imu.Rotations) {
	arm_0_2_rotation, arm_1_3_rotation := transformRollPitch(rotattions.Roll, rotattions.Pitch)

	fc.arm_0_2_PID.SetSetPoint(arm_0_2_rotation)
	fc.arm_1_3_PID.SetSetPoint(arm_1_3_rotation)
}

func (fc *FlightControl) SetThrottle(throttle float64) {
	fc.throttle = fc.throttle*(1-fc.throttleLowPassFilter) + fc.throttleLowPassFilter*throttle
}

func (fc *FlightControl) pidMotorsPowers() {
	fc.escs.SetThrottles([]float64{fc.pidThrottles[0], fc.pidThrottles[1], fc.pidThrottles[2], fc.pidThrottles[3]})
}

func (fc *FlightControl) SetMotorsPowers() {
	if time.Since(fc.motorsArmingTime) >= 0 {
		fc.pidMotorsPowers()
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

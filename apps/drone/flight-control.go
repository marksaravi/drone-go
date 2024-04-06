package drone

import (
	"time"

	"github.com/marksaravi/drone-go/devices/imu"
)

const (
	MOTORS_OFF         = false
	MOTORS_ONN         = true
	THROTTLE_HYSTERSYS = 0.5
)

type FlightControl struct {
	pid              *PidController
	escs             escs
	motorsArmingTime time.Time
}

func NewFlightControl(escs escs, minFlightThrottle, maxThrottle float64, pidConfigs PIDConfigs) *FlightControl {
	fc := &FlightControl{
		pid:  NewPIDControl(pidConfigs, minFlightThrottle, maxThrottle),
		escs: escs,
	}

	fc.turnOnMotors(false)
	return fc
}

func (fc *FlightControl) SetRotations(rotattions imu.Rotations) {
	fc.pid.SetRotations(rotattions)
	fc.pid.CalcESCThrottles()
}

func (fc *FlightControl) SetTargetRotations(rotattions imu.Rotations) {
	fc.pid.SetTargetRotations(rotattions)
}

func (fc *FlightControl) SetThrottle(throttle float64) {
	fc.pid.SetThrottle(throttle)
}

func (fc *FlightControl) ApplyThrottlesToESCs() {
	if time.Since(fc.motorsArmingTime) >= 0 {
		fc.escs.SetThrottles(fc.pid.GetThrottles())
	} else {
		fc.escs.SetThrottles([]float64{0, 0, 0, 0})
	}
}

func (fc *FlightControl) turnOnMotors(motorsOn bool) {
	if motorsOn && fc.pid.throttle < 2 {
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

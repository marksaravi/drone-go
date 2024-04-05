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
	pid               *PIDControl
	escs              escs
	throttle          float64
	minFlightThrottle float64
	motorsArmingTime  time.Time
}

func NewFlightControl(escs escs, minFlightThrottle float64, pidConfigs PIDConfigs) *FlightControl {
	fc := &FlightControl{
		pid:               NewPIDControl(pidConfigs),
		minFlightThrottle: minFlightThrottle,
		escs:              escs,
		throttle:          0,
	}

	fc.turnOnMotors(false)
	return fc
}

func (fc *FlightControl) SetRotations(rotattions imu.Rotations) {
	fc.pid.SetRotations(rotattions)
}

func (fc *FlightControl) SetTargetRotations(rotattions imu.Rotations) {
	fc.pid.SetTargetRotations(rotattions)
}

func (fc *FlightControl) SetThrottle(throttle float64) {
	fc.throttle = throttle
	fc.pid.SetThrottle(throttle)
	if time.Since(fc.motorsArmingTime) < 0 {
		fc.escs.SetThrottles([]float64{0, 0, 0, 0})
		return
	}

	if fc.throttle < fc.minFlightThrottle {
		fc.pid.ResetI()
		fc.escs.SetThrottles([]float64{fc.throttle, fc.throttle, fc.throttle, fc.throttle})
	} else {
		fc.escs.SetThrottles(fc.pid.CalcThrottles())
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

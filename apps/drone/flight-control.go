package drone

import (
	"fmt"
	"time"

	"github.com/marksaravi/drone-go/devices/imu"
)

const (
	MOTORS_OFF         = false
	MOTORS_ONN         = true
	THROTTLE_HYSTERSYS = 0.5
)

type FlightState interface {
	SetRotations(rotattions imu.Rotations)
	SetTargetRotations(rotattions imu.Rotations)
	SetThrottle(throttle float64)
	ApplyESCThrottles()
	Reset(params map[string]bool)
}

type FlightControl struct {
	escs     escs
	throttle float64
	// zeroThrottleState   FlightState
	// lowThrottleState    FlightState
	// flightThrottleState FlightState
	// flightState         FlightState
	// minFlightThrottle   float64
	motorsArmingTime time.Time
}

func NewFlightControl(escs escs, minFlightThrottle float64, pidConfigs PIDConfigs) *FlightControl {
	fc := &FlightControl{
		// minFlightThrottle: minFlightThrottle,
		escs:     escs,
		throttle: 0,
	}

	// fc.zeroThrottleState = &ZeroThrottleState{
	// 	safeZeroStart: false,
	// 	flightControl: fc,
	// }

	// fc.lowThrottleState = &LowThrottleState{
	// 	flightControl: fc,
	// }

	// fc.flightThrottleState = &FlightThrottleState{
	// 	flightControl: fc,
	// 	pid:           NewPIDControl(pidConfigs),
	// }
	fc.turnOnMotors(false)
	return fc
}

// func (fc *FlightControl) SetState(fs FlightState, throttle float64) {
// 	fc.flightState = fs
// 	fc.flightState.Reset(nil)
// 	fc.flightState.SetThrottle(throttle)
// }

func (fc *FlightControl) SetRotations(rotattions imu.Rotations) {
	// fc.SetRotations(rotattions)
}

func (fc *FlightControl) SetTargetRotations(rotattions imu.Rotations) {
	// fc.flightState.SetTargetRotations(rotattions)
}

func (fc *FlightControl) SetThrottle(throttle float64) {
	fc.throttle = throttle
	if time.Since(fc.motorsArmingTime) < 0 {
		fmt.Println("zero throttle")
		fc.escs.SetThrottles([]float64{0, 0, 0, 0})
		return
	}

	fc.escs.SetThrottles([]float64{fc.throttle, fc.throttle, fc.throttle, fc.throttle})
	fmt.Println(fc.throttle)
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

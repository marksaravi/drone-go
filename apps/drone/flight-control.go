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

type FlightState interface {
	SetRotations(rotattions imu.Rotations)
	SetTargetRotations(rotattions imu.Rotations)
	SetThrottle(throttle float64)
	ApplyESCThrottles()
	Reset(params map[string]bool)
}

type FlightControl struct {
	escs                escs
	zeroThrottleState   FlightState
	lowThrottleState    FlightState
	flightThrottleState FlightState
	flightState         FlightState
	minFlightThrottle   float64
	motorsOnTime        time.Time
}

func NewFlightControl(escs escs, minFlightThrottle float64, pidConfigs PIDConfigs) *FlightControl {
	fc := &FlightControl{
		minFlightThrottle: minFlightThrottle,
		escs:              escs,
	}

	fc.zeroThrottleState = &ZeroThrottleState{
		safeZeroStart: false,
		flightControl: fc,
	}

	fc.lowThrottleState = &LowThrottleState{
		flightControl: fc,
	}

	fc.flightThrottleState = &FlightThrottleState{
		flightControl: fc,
		pid:           NewPIDControl(pidConfigs),
	}

	fc.SetToZeroThrottleState(MOTORS_OFF)
	return fc
}

func (fc *FlightControl) SetState(fs FlightState, throttle float64) {
	fc.flightState = fs
	fc.flightState.Reset(nil)
	fc.flightState.SetThrottle(throttle)
}

func (fc *FlightControl) SetRotations(rotattions imu.Rotations) {
	fc.flightState.SetRotations(rotattions)
}

func (fc *FlightControl) SetTargetRotations(rotattions imu.Rotations) {
	fc.flightState.SetTargetRotations(rotattions)
}

func (fc *FlightControl) SetThrottle(throttle float64) {
	fc.flightState.SetThrottle(throttle)
}

func (fc *FlightControl) ApplyESCThrottles() {
	fc.flightState.ApplyESCThrottles()
}

func (fc *FlightControl) SetESCThrottles(throttles []float64) {
	fc.escs.SetThrottles(throttles)
}

func (fc *FlightControl) SetToZeroThrottleState(motorsOn bool) {
	fc.flightState = fc.zeroThrottleState
	fc.flightState.Reset(map[string]bool{"motors-on": motorsOn})
	if motorsOn {
		fc.escs.On()
		fc.motorsOnTime = time.Now()
	} else {
		fc.escs.Off()
	}
}

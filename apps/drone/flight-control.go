package drone

import "fmt"

type FlightState interface {
	SetThrottle(throttle float64)
	ConnectThrottle()
	DisconnectThrottle()
	ResetState()
	ShowState()
}

type FlightControl struct {
	escs                    escs
	noThrottleState         FlightState
	lowThrottleState        FlightState
	flightThrottleState     FlightState
	flightState             FlightState
	flightThrottleThreshold float64
	lowThrottleThreshold    float64
}

func NewFlightControl(escs escs, minFlightThrottle float64) *FlightControl {
	const HYSTERSYS_GAP = 0.5
	fc := &FlightControl{
		flightThrottleThreshold: minFlightThrottle + HYSTERSYS_GAP,
		lowThrottleThreshold:    minFlightThrottle - HYSTERSYS_GAP,
		escs:                    escs,
	}
	fmt.Println("Flight Control: ", fc.flightThrottleThreshold, fc.lowThrottleThreshold)

	fc.noThrottleState = &NoThrottleState{
		flightControl: fc,
	}

	fc.lowThrottleState = &LowThrottleState{
		safeZeroStart: false,
		flightControl: fc,
	}

	fc.flightThrottleState = &FlightThrottleState{
		flightControl: fc,
	}

	fc.SetState(fc.noThrottleState)
	return fc
}

func (fc *FlightControl) SetState(fs FlightState) {
	fs.ResetState()
	fc.flightState = fs
	fc.flightState.ShowState()
}

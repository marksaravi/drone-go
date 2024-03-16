package drone

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

func NewFlightControl(escs escs) *FlightControl {
	fc := &FlightControl{
		escs: escs,
	}

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

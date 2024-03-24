package drone

const (
	MOTORS_OFF         = false
	MOTORS_ONN         = true
	THROTTLE_HYSTERSYS = 0.5
)

type FlightState interface {
	SetThrottle(throttle float64)
	Reset(params map[string]bool)
}

type FlightControl struct {
	escs                escs
	zeroThrottleState   FlightState
	lowThrottleState    FlightState
	flightThrottleState FlightState
	flightState         FlightState
	minFlightThrottle   float64
}

func NewFlightControl(escs escs, minFlightThrottle float64) *FlightControl {
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
	}

	fc.SetToZeroThrottleState(MOTORS_OFF)
	return fc
}

func (fc *FlightControl) SetState(fs FlightState, throttle float64) {
	fc.flightState = fs
	fc.flightState.Reset(nil)
	fc.flightState.SetThrottle(throttle)
}

func (fc *FlightControl) SetThrottles(throttle float64) {
	fc.escs.SetThrottles([]float64{throttle, throttle, throttle, throttle})
}

func (fc *FlightControl) SetToZeroThrottleState(motorsOn bool) {
	fc.flightState = fc.zeroThrottleState
	fc.flightState.Reset(map[string]bool{"motors-on": motorsOn})
	if motorsOn {
		fc.escs.On()
	} else {
		fc.escs.Off()
	}
}

package drone

type FlightState interface {
	SetThrottle(throttle float64)
}

type FlightControl struct {
	noThrottleState     FlightState
	lowThrottleState    FlightState
	flightThrottleState FlightState
	flightState         FlightState
}

func NewFlightControl() *FlightControl {
	fc := &FlightControl{}
	fc.noThrottleState = &NoThrottleState{
		flightControl: fc,
	}
	fc.lowThrottleState = &LowThrottleState{
		flightControl: fc,
	}
	fc.flightThrottleState = &FlightThrottleState{
		flightControl: fc,
	}
	fc.flightState = fc.noThrottleState
	return fc
}

type NoThrottleState struct {
	flightControl *FlightControl
}

func (fs *NoThrottleState) SetThrottle(throttle float64) {

}

type LowThrottleState struct {
	flightControl *FlightControl
}

func (fs *LowThrottleState) SetThrottle(throttle float64) {

}

type FlightThrottleState struct {
	flightControl *FlightControl
}

func (fs *FlightThrottleState) SetThrottle(throttle float64) {

}

package drone

import "fmt"

type FlightThrottleState struct {
	flightControl *FlightControl
	prevThrottle  float64
}

func (fs *FlightThrottleState) Reset(params map[string]bool) {
	fs.prevThrottle = fs.flightControl.minFlightThrottle
	fmt.Println("FLIGHT THROTTLE STATE")
}

func (fs *FlightThrottleState) SetThrottle(throttle float64) {
	if throttle < fs.flightControl.minFlightThrottle-THROTTLE_HYSTERSYS {
		fs.flightControl.SetState(fs.flightControl.lowThrottleState, throttle)
		return
	}
	const K = float64(0.45)
	filteredThrottle := throttle*K + fs.prevThrottle*(1-K)
	fs.prevThrottle = filteredThrottle
	fs.flightControl.SetThrottles(filteredThrottle)
}

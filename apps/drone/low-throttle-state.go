package drone

import "fmt"

type LowThrottleState struct {
	flightControl *FlightControl
}

func (fs *LowThrottleState) Reset(params map[string]bool) {
	fmt.Println("LOW THROTTLE STATE")
}

func (fs *LowThrottleState) SetThrottle(throttle float64) {
	if throttle > fs.flightControl.minFlightThrottle {
		fs.flightControl.SetState(fs.flightControl.flightThrottleState, throttle)
	} else {
		fs.flightControl.SetThrottles(throttle)
	}
}

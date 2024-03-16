package drone

import "log"

type FlightThrottleState struct {
	flightControl *FlightControl
}

func (fs *FlightThrottleState) ShowState() {
	log.Println("Flight Throttle State")
}

func (fs *FlightThrottleState) ResetState() {
}

func (fs *FlightThrottleState) SetThrottle(throttle float64) {
	if throttle < fs.flightControl.lowThrottleThreshold {
		fs.flightControl.SetState(fs.flightControl.lowThrottleState)
		return
	}
	fs.flightControl.escs.SetThrottles([]float64{throttle, throttle, throttle, throttle})
}
func (fs *FlightThrottleState) ConnectThrottle() {
}

func (fs *FlightThrottleState) DisconnectThrottle() {
	log.Println("FLIGHT THROTTLE DISCONNECT")
	fs.flightControl.SetState(fs.flightControl.noThrottleState)
}

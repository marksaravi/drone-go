package drone

import "log"

type FlightThrottleState struct {
	prevThrottle  float64
	flightControl *FlightControl
}

func (fs *FlightThrottleState) ShowState() {
	log.Println("Flight Throttle State")
}

func (fs *FlightThrottleState) ResetState() {
	fs.prevThrottle = fs.flightControl.flightThrottleThreshold
}

func (fs *FlightThrottleState) SetThrottle(throttle float64) {
	const K = float64(0.45)
	filteredThrottle := throttle*K + fs.prevThrottle*(1-K)
	fs.prevThrottle = filteredThrottle
	if filteredThrottle < fs.flightControl.lowThrottleThreshold {
		fs.flightControl.SetState(fs.flightControl.lowThrottleState)
		return
	}
	fs.flightControl.escs.SetThrottles([]float64{filteredThrottle, filteredThrottle, filteredThrottle, filteredThrottle})
}

func (fs *FlightThrottleState) ConnectThrottle() {
}

func (fs *FlightThrottleState) DisconnectThrottle() {
	log.Println("FLIGHT THROTTLE DISCONNECT")
	fs.flightControl.SetState(fs.flightControl.noThrottleState)
}

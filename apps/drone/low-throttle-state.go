package drone

import "log"

type LowThrottleState struct {
	flightControl *FlightControl
	safeZeroStart bool
}

func (fs *LowThrottleState) ShowState() {
	log.Println("Low Throttle State")
}

func (fs *LowThrottleState) ResetState() {
	fs.safeZeroStart = false
	fs.flightControl.escs.On()
}

func (fs *LowThrottleState) SetThrottle(throttle float64) {
	if throttle > fs.flightControl.flightThrottleThreshold {
		fs.flightControl.SetState(fs.flightControl.flightThrottleState)
	}
}

func (fs *LowThrottleState) ConnectThrottle() {
}

func (fs *LowThrottleState) DisconnectThrottle() {
	fs.flightControl.SetState(fs.flightControl.noThrottleState)
}

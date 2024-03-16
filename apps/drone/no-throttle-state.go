package drone

import "log"

type NoThrottleState struct {
	flightControl *FlightControl
}

func (fs *NoThrottleState) ShowState() {
	log.Println("No Throttle State")
}

func (fs *NoThrottleState) ResetState() {
	fs.flightControl.escs.Off()
}

func (fs *NoThrottleState) SetThrottle(throttle float64) {

}

func (fs *NoThrottleState) ConnectThrottle() {
	fs.flightControl.SetState(fs.flightControl.lowThrottleState)
}

func (fs *NoThrottleState) DisconnectThrottle() {
	fs.flightControl.escs.Off()
}

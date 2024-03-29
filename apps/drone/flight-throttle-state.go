package drone

import (
	"fmt"

	"github.com/marksaravi/drone-go/devices/imu"
)

type FlightThrottleState struct {
	flightControl *FlightControl
	prevThrottle  float64
	pid           *PIDControl
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
	fs.pid.SetThrottle(throttle)
}

func (fs *FlightThrottleState) SetRotations(rotattions imu.Rotations) {
	fs.pid.SetRotations(rotattions)
}
func (fs *FlightThrottleState) SetTargetRotations(rotattions imu.Rotations) {
	fs.pid.SetTargetRotations(rotattions)
}
func (fs *FlightThrottleState) ApplyESCThrottles() {
	fs.flightControl.SetESCThrottles(fs.pid.CalcThrottles())
}

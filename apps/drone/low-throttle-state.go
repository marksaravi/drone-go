package drone

import (
	"fmt"

	"github.com/marksaravi/drone-go/devices/imu"
	"github.com/marksaravi/drone-go/utils"
)

type LowThrottleState struct {
	flightControl *FlightControl
	throttle      float64
}

func (fs *LowThrottleState) Reset(params map[string]bool) {
	fmt.Println("LOW THROTTLE STATE")
}

var lowThrottleDisplay = utils.WithDataPerSecond(5)

func (fs *LowThrottleState) SetThrottle(throttle float64) {
	// fs.throttle = throttle
	// if throttle > fs.flightControl.minFlightThrottle {
	// 	fs.flightControl.SetState(fs.flightControl.flightThrottleState, throttle)
	// 	return
	// }
	// if lowThrottleDisplay.IsTime() {
	// 	fmt.Println(throttle)
	// }
}

func (fs *LowThrottleState) SetRotations(rotattions imu.Rotations) {}

func (fs *LowThrottleState) SetTargetRotations(rotattions imu.Rotations) {}

func (fs *LowThrottleState) ApplyESCThrottles() {
	// fs.flightControl.SetESCThrottles([]float64{fs.throttle, fs.throttle, fs.throttle, fs.throttle})
}

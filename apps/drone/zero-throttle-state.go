package drone

import (
	"fmt"
	"time"
)

type ZeroThrottleState struct {
	flightControl *FlightControl
	safeZeroStart bool
	motorsOn      bool
}

func (fs *ZeroThrottleState) Reset(params map[string]bool) {
	fs.safeZeroStart = false
	fs.motorsOn = params["motors-on"]
	fmt.Println("ZERO THROTTLE STATE MOTORS ON ", fs.motorsOn)
}

func (fs *ZeroThrottleState) SetThrottle(throttle float64) {
	if !fs.safeZeroStart && throttle == 0 && fs.motorsOn {
		fs.setFlightSafe()
	} else if fs.safeZeroStart && throttle > 0 && throttle <= 3 {
		fs.flightControl.SetState(fs.flightControl.lowThrottleState, throttle)
	}
}

func (fs *ZeroThrottleState) setFlightSafe() {
	for time.Since(fs.flightControl.motorsOnTime) < time.Second*3 {
		time.Sleep(time.Millisecond * 100)
	}
	fs.safeZeroStart = true
	fmt.Println("SAFE THROTTLE")
}

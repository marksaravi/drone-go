package drone

import "fmt"

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
		fs.safeZeroStart = true
		fmt.Println("SAFE THROTTLE")
	} else if fs.safeZeroStart && throttle > 1 {
		fs.flightControl.SetState(fs.flightControl.lowThrottleState, throttle)
	}
}

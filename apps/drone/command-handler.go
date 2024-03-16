package drone

import (
	"fmt"

	"github.com/marksaravi/drone-go/constants"
)

func (d *droneApp) applyCommands(commands []byte) {
	if d.offMotors(commands) {
		return
	}
	if d.onMotors(commands) {
		return
	}
	d.setCommands(commands)
}

func (d *droneApp) onMotors(commands []byte) bool {
	if commands[5] == 1 {
		d.flightControl.flightState.ConnectThrottle()
		return true
	}
	return false
}

func (d *droneApp) offMotors(commands []byte) bool {
	if commands[5] == 16 {
		d.flightControl.flightState.DisconnectThrottle()
		return true
	}
	return false
}

func (d *droneApp) setCommands(commands []byte) {
	throttle := float64(commands[3]) / float64(constants.THROTTLE_RAW_READING_MAX) * d.maxApplicableThrottle
	fmt.Println(throttle)
}

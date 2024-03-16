package drone

import (
	"fmt"

	"github.com/marksaravi/drone-go/apps/common"
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
	if commands[9] == 1 {
		d.flightControl.flightState.ConnectThrottle()
		return true
	}
	return false
}

func (d *droneApp) offMotors(commands []byte) bool {
	if commands[9] == 16 {
		d.flightControl.flightState.DisconnectThrottle()
		return true
	}
	return false
}

func (d *droneApp) setCommands(commands []byte) {
	roll := common.CalcRotationFromRawJoyStickRaw(commands[0:2], d.rollMidValue, d.rotationRange)
	pitch := common.CalcRotationFromRawJoyStickRaw(commands[2:4], d.pitchlMidValue, d.rotationRange)
	yaw := common.CalcRotationFromRawJoyStickRaw(commands[4:6], d.yawMidValue, d.rotationRange)
	throttle := common.CalcThrottleFromRawJoyStickRaw(commands[6:8], d.maxThrottle)
	fmt.Printf("%6.2f, %6.2f, %6.2f, %6.2f \n", roll, pitch, yaw, throttle)
}

func RawJoystickRotation(commands []byte) int {
	return int(uint16(commands[1]) | (uint16(commands[0]) << 8))
}

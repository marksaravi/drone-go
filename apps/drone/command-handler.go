package drone

import (
	"fmt"
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
	rawRoll, rawPitch, rawYaw, rawThrottle := RawRollPitchYawThrottle(commands)
	fmt.Printf("%4d, %4d, %4d, %4d \n", rawRoll, rawPitch, rawYaw, rawThrottle)
}

func RawRollPitchYawThrottle(commands []byte) (rawRoll, rawPitch, rawYaw, rawThrottle uint16) {
	rawRoll = uint16(commands[1]) | (uint16(commands[0]) << 8)
	rawPitch = uint16(commands[3]) | (uint16(commands[2]) << 8)
	rawYaw = uint16(commands[5]) | (uint16(commands[4]) << 8)
	rawThrottle = uint16(commands[7]) | (uint16(commands[6]) << 8)
	return
}

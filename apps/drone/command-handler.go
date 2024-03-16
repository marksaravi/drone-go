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
	roll := CalcRotationFromRawJoyStickRaw(commands[0:2], d.rollMidValue, d.rotationRange)
	pitch := CalcRotationFromRawJoyStickRaw(commands[2:4], d.pitchlMidValue, d.rotationRange)
	yaw := CalcRotationFromRawJoyStickRaw(commands[4:6], d.yawMidValue, d.rotationRange)
	throttle := CalcThrottleFromRawJoyStickRaw(commands[6:8], d.maxThrottle)
	fmt.Printf("%6.2f, %6.2f, %6.2f, %6.2f \n", roll, pitch, yaw, throttle)
}

func RawJoystickRotation(commands []byte) int {
	return int(uint16(commands[1]) | (uint16(commands[0]) << 8))
}

func CalcRotationFromRawJoyStickRaw(commands []byte, midValue int, rotationRange float64) float64 {
	rawValue := RawJoystickRotation(commands)
	rawValue -= midValue
	if rawValue < -midValue {
		rawValue = -midValue
	}
	if rawValue > midValue {
		rawValue = midValue
	}
	return float64(rawValue) / float64(midValue) * rotationRange
}
func CalcThrottleFromRawJoyStickRaw(commands []byte, maxThrottle float64) float64 {
	rawValue := RawJoystickRotation(commands)
	return float64(rawValue) / float64(constants.JOYSTICK_DIGITAL_MAX_VALUE) * maxThrottle
}

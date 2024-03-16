package common

import (
	"github.com/marksaravi/drone-go/constants"
)

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

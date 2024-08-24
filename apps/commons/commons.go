package commons

import (
	"github.com/marksaravi/drone-go/constants"
)

func Uint16ToBytes(x uint16) (low, high byte) {
	low = byte(x)
	high = byte(x >> 8)
	return
}

func RawJoystickRotation(commands []byte) int {
	return int(uint16(commands[0]) | (uint16(commands[1]) << 8))
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
	return float64(rawValue) / float64(constants.JOY_STICK_INPUT_RANGE) * maxThrottle
}

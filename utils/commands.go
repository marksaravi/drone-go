package utils

import "github.com/marksaravi/drone-go/constants"

func desrializeJoyStickCommand(hByte, lByte byte) uint16 {
	h := uint16(hByte) << 8
	l := uint16(lByte)
	return h | l
}

func CommandToRotation(hByte, lByte byte, maxRotation float64) float64 {
	halfRange := float64(constants.JOY_STICK_INPUT_RANGE / 2)
	v := float64(desrializeJoyStickCommand(hByte, lByte)) - halfRange
	return v * maxRotation / halfRange
}

func CommandToThrottle(hByte, lByte byte, maxThrottle float64) float64 {
	fullRange := float64(constants.JOY_STICK_INPUT_RANGE)
	v := float64(desrializeJoyStickCommand(hByte, lByte))
	return v * maxThrottle / fullRange
}

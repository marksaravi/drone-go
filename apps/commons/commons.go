package commons

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


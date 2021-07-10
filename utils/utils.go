package utils

// TowsComplementUint8ToInt16 converts 2's complement H and L uint8 to signed int16
func TowsComplementUint8ToInt16(h, l uint8) int16 {
	var h16 uint16 = uint16(h)
	var l16 uint16 = uint16(l)

	return int16((h16 << 8) | l16)
}

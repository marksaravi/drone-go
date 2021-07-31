package utils

import "unsafe"

// TowsComplementUint8ToInt16 converts 2's complement H and L uint8 to signed int16
func TowsComplementUint8ToInt16(h, l uint8) int16 {
	var h16 uint16 = uint16(h)
	var l16 uint16 = uint16(l)

	return int16((h16 << 8) | l16)
}

func Min(a byte, b byte) byte {
	if a <= b {
		return a
	}
	return b
}

func Max(a byte, b byte) byte {
	if a >= b {
		return a
	}
	return b
}

func FloatArrayToByteArray(floatArray []float32) []byte {
	faLen := len(floatArray)
	byteArray := make([]byte, faLen*4)
	for i := 0; i < faLen; i++ {
		ba := int32ToByteArray(float32ToInt32(floatArray[i]))
		for j := 0; j < 4; j++ {
			byteArray[i*4+j] = ba[j]
		}
	}
	return byteArray
}

func ByteArrayToFloat32Array(byteArray []byte) []float32 {
	baLen := len(byteArray)
	floatArray := make([]float32, baLen/4)
	for i := 0; i < baLen; i += 4 {
		floatArray[i/4] = int32ToFloat32(byteArrayToInt32(byteArray[i : i+4]))
	}
	return floatArray
}

func float32ToInt32(f float32) int32 {
	type pi32 = *int32
	var pi pi32 = pi32(unsafe.Pointer(&f))
	return *pi
}

func int32ToFloat32(i int32) float32 {
	type pf32 = *float32
	var pf pf32 = pf32(unsafe.Pointer(&i))
	return *pf
}

func int32ToByteArray(i int32) []byte {
	ba := make([]byte, 4)
	const mask = 0b00000000000000000000000011111111
	ba[0] = byte(i & mask)
	ba[1] = byte((i >> 8) & mask)
	ba[2] = byte((i >> 16) & mask)
	ba[3] = byte((i >> 24) & mask)
	return ba
}

func byteArrayToInt32(ba []byte) int32 {
	var i int32 = 0
	i = i | int32(ba[0])
	i = i | (int32(ba[1]) << 8)
	i = i | (int32(ba[2]) << 16)
	i = i | (int32(ba[3]) << 24)
	return i
}

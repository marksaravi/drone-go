package utils

import (
	"bytes"
	"encoding/binary"
	"math"
)

func UInt16ToBytes(i uint16) []byte {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, i)
	return buf
}

func UInt64ToBytes(i uint64) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, i)
	return buf
}

func UInt64FromBytes(bytes []byte) uint64 {
	return binary.LittleEndian.Uint64(bytes)
}

func UInt32ToBytes(i uint32) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, i)
	return buf
}

func UInt32FromBytes(bytes []byte) uint32 {
	return binary.LittleEndian.Uint32(bytes)
}

func Float32ToBytes(f float32) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, f)
	return buf.Bytes()
}

func Float32FromBytes(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}

func BoolArrayToByte(bools [8]bool) byte {
	var res byte = 0
	for i := 0; i < 8; i++ {
		if !bools[i] {
			continue
		}
		var mask byte = 1
		mask <<= i
		res |= mask
	}
	return res
}

func BoolArrayFromByte(b byte) [8]bool {
	var res [8]bool = [8]bool{}
	for i := 0; i < 8; i++ {
		var mask byte = 1 << i
		res[i] = (mask & b) > 0
	}
	return res
}

func Float64ToRoundedFloat32Bytes(x float64) []byte {
	var v float32 = float32(math.Round(x*100) / 100)
	return Float32ToBytes(v)
}
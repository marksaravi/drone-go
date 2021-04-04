package utils

import (
	"encoding/binary"
	"math"
	"unsafe"

	"github.com/MarkSaravi/drone-go/types"
)

// TowsComplementToInt converts 16 bit 2's complement to signed int
func TowsComplementToInt(a uint16) int16 {
	return *(*int16)(unsafe.Pointer(&a))
}

// IntToTowsComplement converts 16 bit signed int to 2's complement
func IntToTowsComplement(a int16) uint16 {
	return *(*uint16)(unsafe.Pointer(&a))
}

// TowsComplementBytesToInt converts 2's complement H and L uint8 to signed int16
func TowsComplementBytesToInt(h, l uint8) int16 {
	b := binary.BigEndian.Uint16([]uint8{h, l})
	return TowsComplementToInt(b)
}

// IntToTowsComplementBytes converts 2's complement H and L uint8 to signed int16
func IntToTowsComplementBytes(a int16) (h, l uint8) {
	b := make([]uint8, 2)
	binary.LittleEndian.PutUint16(b[0:], IntToTowsComplement(a))
	h = b[1]
	l = b[0]
	return
}

// CalcVectorLen calculate the length of a vector
func CalcVectorLen(v types.XYZ) float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

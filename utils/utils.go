package utils

import (
	"encoding/binary"
	"math"
	"unsafe"

	"github.com/MarkSaravi/drone-go/types/sensore"
)

// TowsComplementToInt converts 16 bit 2's complement to signed int
func TowsComplementToInt(a uint16) int16 {
	return *(*int16)(unsafe.Pointer(&a))
}

// IntToTowsComplement converts 16 bit signed int to 2's complement
func IntToTowsComplement(a int16) uint16 {
	return *(*uint16)(unsafe.Pointer(&a))
}

// TowsComplementBytesToInt converts 2's complement H and L byte to signed int16
func TowsComplementBytesToInt(h, l byte) int16 {
	b := binary.BigEndian.Uint16([]byte{h, l})
	return TowsComplementToInt(b)
}

// IntToTowsComplementBytes converts 2's complement H and L byte to signed int16
func IntToTowsComplementBytes(a int16) (h, l byte) {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b[0:], IntToTowsComplement(a))
	h = b[1]
	l = b[0]
	return
}

// CalcVectorLen calculate the length of a vector
func CalcVectorLen(v sensore.Data) float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

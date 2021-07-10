package utils

import (
	"encoding/binary"
	"fmt"
	"math"
	"time"
	"unsafe"

	"github.com/MarkSaravi/drone-go/types"
)

var lastReading = time.Now()

func Print(v []float64, msInterval int) {
	if time.Since(lastReading) >= time.Millisecond*time.Duration(msInterval) {
		fmt.Println(v)
		lastReading = time.Now()
	}
}

// IntToTowsComplement converts 16 bit signed int to 2's complement
func IntToTowsComplement(a int16) uint16 {
	return *(*uint16)(unsafe.Pointer(&a))
}

// TowsComplementUint8ToInt16 converts 2's complement H and L uint8 to signed int16
func TowsComplementUint8ToInt16(h, l uint8) int16 {
	var h16 uint16 = uint16(h)
	var l16 uint16 = uint16(l)

	return int16((h16 << 8) | l16)
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

package utils

import (
	"bytes"
	"encoding/binary"
	"math"

	"github.com/marksaravi/drone-go/models"
)

func UInt64ToBytes(i uint64) [8]byte {
	buf := [8]byte{0, 0, 0, 0, 0, 0, 0, 0}
	binary.LittleEndian.PutUint64(buf[:], i)
	return buf
}

func UInt64FromBytes(bytes [8]byte) uint64 {
	return binary.LittleEndian.Uint64(bytes[:])
}

func Int64ToBytes(i int64) [8]byte {
	return UInt64ToBytes(uint64(i))
}

func Int64FromBytes(bytes [8]byte) int64 {
	return int64(UInt64FromBytes(bytes))
}

func UInt32ToBytes(i uint32) [4]byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, i)
	return SliceToArray4(buf)
}

func UInt32FromBytes(bytes [4]byte) uint32 {
	return binary.LittleEndian.Uint32(bytes[:])
}

func Float32ToBytes(f float32) [4]byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, f)
	return SliceToArray4(buf.Bytes())
}

func Float32FromBytes(bytes [4]byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes[:])
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

func Float64ToRoundedFloat32Bytes(x float64) [4]byte {
	var v float32 = float32(math.Round(x*100) / 100)
	return Float32ToBytes(v)
}

func SliceToArray4(slice []byte) [4]byte {
	array := [4]byte{}
	copy(array[:], slice)
	return array
}

func SliceToArray8(slice []byte) [8]byte {
	array := [8]byte{}
	copy(array[:], slice)
	return array
}

func SliceToArray32(slice []byte) models.Payload {
	array := models.Payload{}
	copy(array[:], slice)
	return array
}

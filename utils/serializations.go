package utils

import (
	"encoding/binary"
	"math"
	"time"
)

func SerializeFloat64(f float64) []byte {
	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, uint16(int16(math.Round(f*10))))
	return data
}

func DeSerializeFloat64(data []byte) float64 {
	return float64(int16(binary.LittleEndian.Uint16(data))) / 10
}

func SerializeInt(n int16) []byte {
	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, uint16(n))
	return data
}

func DeSerializeInt(data []byte) int16 {
	return int16(binary.LittleEndian.Uint16(data))
}

func SerializeDuration(dur time.Duration) []byte {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, uint32(dur.Microseconds()/100))
	return data
}

func DeSerializeDuration(data []byte) time.Duration {
	return time.Duration(int32(binary.LittleEndian.Uint32(data))) * time.Microsecond * 100
}

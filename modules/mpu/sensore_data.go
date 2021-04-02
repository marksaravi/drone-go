package mpu

import "github.com/MarkSaravi/drone-go/types"

type SensorData struct {
	numOfData uint32
	data      types.XYZ
	buffer    []types.XYZ
}

func NewSensorData(bufferSize uint8) SensorData {
	return SensorData{
		numOfData: 0,
		data:      types.XYZ{X: 0, Y: 0, Z: 0},
		buffer:    make([]types.XYZ, bufferSize),
	}
}

func (s *SensorData) PushToFront(xyz types.XYZ) {
	s.buffer = append([]types.XYZ{xyz}, s.buffer[:len(s.buffer)-1]...)
	s.numOfData++
	s.ProcessData()
}

func (s *SensorData) GetBuffer() []types.XYZ {
	return s.buffer
}

func (s *SensorData) ProcessData() {
	if s.numOfData < uint32(BUFFER_LEN) {
		return
	}
}

package mpu

import (
	"github.com/MarkSaravi/drone-go/types"
)

// MPU is interface to mpu mems
type MPU interface {
	Close() error
	InitDevice() error
	Start()
	ReadRawData() ([]byte, error)
	ReadData() (acc types.XYZ, isAccDataReady bool, gyro types.XYZ, isGyroDataReady bool, err error)
	WhoAmI() (string, byte, error)
}

type SensorData struct {
	isValid bool
	data    types.XYZ
	buffer  []types.XYZ
}

type MpuHandler struct {
	mpu        MPU
	acc        types.XYZ
	isAccValid bool
	gyro       types.XYZ
}

func initBuffer(size int) []types.XYZ {
	b := []types.XYZ{}
	for i := 0; i < size; i++ {
		b = append(b, types.XYZ{X: 0, Y: 0, Z: 0})
	}
	return b
}

func New(bufferSize int) SensorData {
	return SensorData{
		isValid: false,
		data:    types.XYZ{X: 0, Y: 0, Z: 0},
		buffer:  initBuffer(bufferSize),
	}
}

func (s *SensorData) PushToFront(xyz types.XYZ) {
	s.buffer = append([]types.XYZ{xyz}, s.buffer[:len(s.buffer)-1]...)
}

func (s *SensorData) GetBuffer() []types.XYZ {
	return s.buffer
}

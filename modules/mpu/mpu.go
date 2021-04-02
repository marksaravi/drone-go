package mpu

import (
	"github.com/MarkSaravi/drone-go/types"
)

const BUFFER_LEN uint8 = 4

// MPU is interface to mpu mems
type MpuDevice interface {
	Close() error
	InitDevice() error
	Start()
	ReadRawData() ([]byte, error)
	ReadData() (acc types.XYZ, isAccDataReady bool, gyro types.XYZ, isGyroDataReady bool, err error)
	WhoAmI() (string, byte, error)
}

type MPU struct {
	dev  MpuDevice
	acc  SensorData
	gyro SensorData
}

func NewMPU(dev MpuDevice) *MPU {
	return &MPU{
		dev:  dev,
		acc:  NewSensorData(BUFFER_LEN),
		gyro: NewSensorData(BUFFER_LEN),
	}
}

func (mpu *MPU) ReadData() {
	accData, isAccDataReady, gyroData, isGyroDataReady, err := mpu.dev.ReadData()
	if err != nil {
		return
	}
	if isAccDataReady {
		mpu.acc.PushToFront(accData)
	}
	if isGyroDataReady {
		mpu.acc.PushToFront(gyroData)
	}
}

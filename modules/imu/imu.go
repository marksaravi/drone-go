package imu

import (
	"github.com/MarkSaravi/drone-go/types"
)

const BUFFER_LEN uint8 = 4

// IMU is interface to imu mems
type ImuDevice interface {
	Close() error
	InitDevice() error
	Start()
	ReadRawData() ([]byte, error)
	ReadData() (acc types.XYZ, isAccDataReady bool, gyro types.XYZ, isGyroDataReady bool, err error)
	WhoAmI() (string, byte, error)
}

type IMU struct {
	Dev         ImuDevice
	acc         SensorData
	gyro        SensorData
	Orientation types.XYZ
}

func NewIMU(dev ImuDevice) *IMU {
	return &IMU{
		Dev:  dev,
		acc:  NewSensorData(BUFFER_LEN),
		gyro: NewSensorData(BUFFER_LEN),
	}
}

func (imu *IMU) ReadData() {
	accData, isAccDataReady, gyroData, isGyroDataReady, err := imu.Dev.ReadData()
	if err != nil {
		return
	}
	if isAccDataReady {
		imu.acc.PushToFront(accData)
	}
	if isGyroDataReady {
		imu.gyro.PushToFront(gyroData)
	}
	imu.Orientation = imu.gyro.data
}

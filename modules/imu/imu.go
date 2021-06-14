package imu

import (
	"github.com/MarkSaravi/drone-go/types"
	"github.com/MarkSaravi/drone-go/utils"
)

const BUFFER_LEN uint8 = 4

// IMU is interface to imu mems
type IMU interface {
	Close()
	GetRotations() (types.ImuRotations, error)
}

type ImuMems interface {
	Close()
	InitDevice() error
	ReadSensorsRawData() ([]byte, error)
	ReadSensors() (types.ImuSensorsData, error)
	WhoAmI() (string, byte, error)
}

type ImuDevice2 struct {
	imuMemes                 ImuMems
	acc                      types.Sensor
	gyro                     types.Sensor
	mag                      types.Sensor
	prevRotations            types.Rotations
	prevGyro                 types.Rotations
	prevReadTime             int64
	readTime                 int64
	lowPassFilterCoefficient float64
}

func NewImuDevice(imuMems ImuMems) ImuDevice2 {
	return ImuDevice2{
		imuMemes: imuMems,
	}
}

func (dev ImuDevice2) Close() {
	dev.imuMemes.Close()
}

func (dev ImuDevice2) GetRotations() (types.ImuRotations, error) {
	imuData, imuError := dev.imuMemes.ReadSensors()
	dg := utils.GyroChanges(imuData)
	gyroRotations := utils.GyroRotations(dg, dev.prevGyro)
	dev.prevGyro = gyroRotations
	accRotations := utils.AccelerometerDataRotations(imuData.Acc.Data)
	prevRotations := dev.prevRotations
	rotations := utils.CalcRotations(
		prevRotations,
		accRotations,
		dg,
		dev.lowPassFilterCoefficient,
	)
	dev.prevRotations = rotations
	return types.ImuRotations{
		PrevRotations: prevRotations,
		Accelerometer: accRotations,
		Gyroscope:     gyroRotations,
		Rotations:     rotations,
		ReadTime:      imuData.ReadTime,
		ReadInterval:  imuData.ReadInterval,
	}, imuError
}

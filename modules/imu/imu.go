package imu

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/types"
	"github.com/MarkSaravi/drone-go/utils"
)

type ImuDevice struct {
	imuMemes                 types.ImuMems
	acc                      types.Sensor
	gyro                     types.Sensor
	mag                      types.Sensor
	prevRotations            types.Rotations
	prevGyro                 types.Rotations
	prevReadTime             int64
	readTime                 int64
	lowPassFilterCoefficient float64
}

func (imudev *ImuDevice) initDeviceReadings(now int64) {
	fmt.Println("initDeviceReadings")
	imudev.prevReadTime = now
	imudev.prevRotations = types.Rotations{Roll: 0, Pitch: 0, Yaw: 0}
	imudev.prevGyro = types.Rotations{Roll: 0, Pitch: 0, Yaw: 0}
}

func (imudev *ImuDevice) readSensors() (types.ImuSensorsData, error) {
	now := time.Now().UnixNano()
	if imudev.readTime < 0 {
		imudev.initDeviceReadings(now)
	} else {
		imudev.prevReadTime = imudev.readTime
	}
	imudev.readTime = now
	acc, gyro, mag, err := imudev.imuMemes.ReadSensors()

	return types.ImuSensorsData{
		Acc:          acc,
		Gyro:         gyro,
		Mag:          mag,
		ReadTime:     imudev.readTime,
		ReadInterval: imudev.readTime - imudev.prevReadTime,
	}, err
}

func NewImuDevice(imuMems types.ImuMems, lowPassFilterCoefficient float64) ImuDevice {
	return ImuDevice{
		imuMemes:                 imuMems,
		prevReadTime:             -1,
		readTime:                 -1,
		lowPassFilterCoefficient: lowPassFilterCoefficient,
	}
}

func (imudev ImuDevice) Close() {
	imudev.imuMemes.Close()
}

func (imudev *ImuDevice) GetRotations() (types.ImuRotations, error) {
	imuData, imuError := imudev.readSensors()
	dg := utils.GyroChanges(imuData)
	gyroRotations := utils.GyroRotations(dg, imudev.prevGyro)
	imudev.prevGyro = gyroRotations
	accRotations := utils.AccelerometerDataRotations(imuData.Acc.Data)
	prevRotations := imudev.prevRotations
	rotations := utils.CalcRotations(
		prevRotations,
		accRotations,
		dg,
		imudev.lowPassFilterCoefficient,
	)
	imudev.prevRotations = rotations
	return types.ImuRotations{
		PrevRotations: prevRotations,
		Accelerometer: accRotations,
		Gyroscope:     gyroRotations,
		Rotations:     rotations,
		ReadTime:      imuData.ReadTime,
		ReadInterval:  imuData.ReadInterval,
	}, imuError
}

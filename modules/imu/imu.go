package imu

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/types"
	"github.com/MarkSaravi/drone-go/utils"
)

type ImuModule struct {
	dev                      types.ImuDevice
	acc                      types.Sensor
	gyro                     types.Sensor
	mag                      types.Sensor
	prevRotations            types.Rotations
	prevGyro                 types.Rotations
	prevReadTime             int64
	readTime                 int64
	lowPassFilterCoefficient float64
}

func (imu *ImuModule) initDeviceReadings(now int64) {
	fmt.Println("initDeviceReadings")
	imu.prevReadTime = now
	imu.prevRotations = types.Rotations{Roll: 0, Pitch: 0, Yaw: 0}
	imu.prevGyro = types.Rotations{Roll: 0, Pitch: 0, Yaw: 0}
}

func (imu *ImuModule) readSensors() (types.ImuSensorsData, error) {
	now := time.Now().UnixNano()
	if imu.readTime < 0 {
		imu.initDeviceReadings(now)
	} else {
		imu.prevReadTime = imu.readTime
	}
	imu.readTime = now
	acc, gyro, mag, err := imu.dev.ReadSensors()

	return types.ImuSensorsData{
		Acc:          acc,
		Gyro:         gyro,
		Mag:          mag,
		ReadTime:     imu.readTime,
		ReadInterval: imu.readTime - imu.prevReadTime,
	}, err
}

func NewIMU(imuMems types.ImuDevice, lowPassFilterCoefficient float64) ImuModule {
	return ImuModule{
		dev:                      imuMems,
		prevReadTime:             -1,
		readTime:                 -1,
		lowPassFilterCoefficient: lowPassFilterCoefficient,
	}
}

func (imu ImuModule) Close() {
	imu.dev.Close()
}

func (imu *ImuModule) GetRotations() (types.ImuRotations, error) {
	imuData, imuError := imu.readSensors()
	dg := utils.GyroChanges(imuData)
	gyroRotations := utils.GyroRotations(dg, imu.prevGyro)
	imu.prevGyro = gyroRotations
	accRotations := utils.AccelerometerDataRotations(imuData.Acc.Data)
	prevRotations := imu.prevRotations
	rotations := utils.CalcRotations(
		prevRotations,
		accRotations,
		dg,
		imu.lowPassFilterCoefficient,
	)
	imu.prevRotations = rotations
	return types.ImuRotations{
		PrevRotations: prevRotations,
		Accelerometer: accRotations,
		Gyroscope:     gyroRotations,
		Rotations:     rotations,
		ReadTime:      imuData.ReadTime,
		ReadInterval:  imuData.ReadInterval,
	}, imuError
}

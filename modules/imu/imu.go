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
	startTime                time.Time
	readTime                 time.Time
	readingInterval          time.Duration
	lowPassFilterCoefficient float64
	readingData              types.ImuReadingQualities
}

func (imu *ImuModule) ResetReadingTimes() {
	imu.startTime = time.Now()
	imu.readTime = imu.startTime
	imu.prevRotations = types.Rotations{Roll: 0, Pitch: 0, Yaw: 0}
	imu.prevGyro = types.Rotations{Roll: 0, Pitch: 0, Yaw: 0}
}

func (imu *ImuModule) updateReadingData(diff time.Duration, err bool) {
	imu.readingData.Total++
	if diff >= imu.readingData.BadIntervalThereshold {
		imu.readingData.BadInterval++
		if diff > imu.readingData.MaxBadInterval {
			imu.readingData.MaxBadInterval = diff
		}
	}
	if err {
		imu.readingData.BadData++
	}
}

func (imu *ImuModule) readSensors() (acc types.SensorData, gyro types.SensorData, mag types.SensorData, err error) {
	acc, gyro, mag, err = imu.dev.ReadSensors()
	return
}

func NewIMU(imuMems types.ImuDevice, config types.FlightConfig) ImuModule {
	readingInterval := time.Duration(int64(time.Second) / int64(config.ImuDataPerSecond))
	badIntervalThereshold := readingInterval + readingInterval/20
	fmt.Println(readingInterval, badIntervalThereshold)
	return ImuModule{
		dev:                      imuMems,
		readTime:                 time.Time{},
		readingInterval:          readingInterval,
		lowPassFilterCoefficient: config.LowPassFilterCoefficient,
		readingData: types.ImuReadingQualities{
			Total:                 0,
			BadInterval:           0,
			MaxBadInterval:        0,
			BadData:               0,
			BadIntervalThereshold: badIntervalThereshold,
		},
	}
}

func (imu ImuModule) Close() {
	imu.dev.Close()
}

func (imu *ImuModule) GetRotations() (bool, types.ImuRotations, error) {
	now := time.Now()
	diff := now.Sub(imu.readTime)
	if diff < imu.readingInterval {
		return false, types.ImuRotations{}, nil
	}
	imu.readTime = now
	acc, gyro, _, imuError := imu.readSensors()
	imu.updateReadingData(diff, imuError != nil)
	dg := utils.GyroChanges(gyro, diff.Nanoseconds())
	gyroRotations := utils.GyroRotations(dg, imu.prevGyro)
	imu.prevGyro = gyroRotations
	accRotations := utils.AccelerometerDataRotations(acc.Data)
	prevRotations := imu.prevRotations
	rotations := utils.CalcRotations(
		prevRotations,
		accRotations,
		dg,
		imu.lowPassFilterCoefficient,
	)
	imu.prevRotations = rotations
	return true, types.ImuRotations{
		PrevRotations: prevRotations,
		Accelerometer: accRotations,
		Gyroscope:     gyroRotations,
		Rotations:     rotations,
		ReadTime:      imu.readTime.UnixNano() - imu.startTime.UnixNano(),
		ReadInterval:  diff.Nanoseconds(),
	}, imuError
}

func (imu *ImuModule) GetReadingQualities() types.ImuReadingQualities {
	return imu.readingData
}

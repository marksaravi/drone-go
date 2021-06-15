package imu

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/types"
	"github.com/MarkSaravi/drone-go/utils"
)

type ImuModule struct {
	dev                      types.ImuDevice
	rotations                types.Rotations
	gyro                     types.Rotations
	startTime                time.Time
	readTime                 time.Time
	readingInterval          time.Duration
	lowPassFilterCoefficient float64
	readingData              types.ImuReadingQualities
}

func (imu *ImuModule) ResetReadingTimes() {
	imu.startTime = time.Now()
	imu.readTime = imu.startTime
	imu.rotations = types.Rotations{Roll: 0, Pitch: 0, Yaw: 0}
	imu.gyro = types.Rotations{Roll: 0, Pitch: 0, Yaw: 0}
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
	imu.gyro = utils.GyroRotations(dg, imu.gyro)
	accRotations := utils.AccelerometerRotations(acc.Data)
	prevRotations := imu.rotations
	imu.rotations = utils.CalcRotations(
		prevRotations,
		accRotations,
		dg,
		imu.lowPassFilterCoefficient,
	)

	return true, types.ImuRotations{
		Accelerometer: accRotations,
		Gyroscope:     imu.gyro,
		Rotations:     imu.rotations,
		ReadTime:      imu.readTime.UnixNano() - imu.startTime.UnixNano(),
		ReadInterval:  diff.Nanoseconds(),
	}, imuError
}

func (imu *ImuModule) GetReadingQualities() types.ImuReadingQualities {
	return imu.readingData
}

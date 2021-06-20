package imu

import (
	"fmt"
	"math"
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
		rotations: types.Rotations{Roll: 0, Pitch: 0, Yaw: 0},
		gyro:      types.Rotations{Roll: 0, Pitch: 0, Yaw: 0},
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
	accData, gyroData, _, devErr := imu.dev.ReadSensors()

	accRotations := AccelerometerRotations(accData.Data)
	dg := GyroChanges(gyroData.Data, diff.Nanoseconds())
	imu.gyro = GyroRotations(dg, imu.gyro)
	rotations := GyroRotations(dg, imu.rotations)
	imu.rotations = rotations

	imu.updateReadingData(diff, devErr != nil)
	return true, types.ImuRotations{
		Accelerometer: accRotations,
		Gyroscope:     imu.gyro,
		Rotations:     imu.rotations,
		ReadTime:      imu.readTime.UnixNano() - imu.startTime.UnixNano(),
		ReadInterval:  diff.Nanoseconds(),
	}, devErr
}

func GyroChanges(gyro types.XYZ, timeInterval int64) types.RotationsChanges {
	dt := goDurToDt(timeInterval)
	return types.RotationsChanges{
		DRoll:  gyro.X * dt,
		DPitch: gyro.Y * dt,
		DYaw:   gyro.Z * dt,
	}
}

func GyroRotations(dGyro types.RotationsChanges, gyro types.Rotations) types.Rotations {
	return types.Rotations{
		Roll:  math.Mod(gyro.Roll+dGyro.DRoll, 360),
		Pitch: math.Mod(gyro.Pitch+dGyro.DPitch, 360),
		Yaw:   math.Mod(gyro.Yaw+dGyro.DYaw, 360),
	}
}

func AccelerometerRotations(data types.XYZ) types.Rotations {
	roll := utils.RadToDeg(math.Atan2(data.Y, data.Z))
	pitch := -utils.RadToDeg(math.Atan2(data.X, data.Z))
	return types.Rotations{
		Roll:  roll,
		Pitch: pitch,
		Yaw:   0,
	}
}

func goDurToDt(d int64) float64 {
	return float64(d) / 1e9
}

func (imu *ImuModule) GetReadingQualities() types.ImuReadingQualities {
	return imu.readingData
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

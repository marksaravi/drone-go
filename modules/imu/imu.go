package imu

import (
	"fmt"
	"math"
	"time"

	"github.com/MarkSaravi/drone-go/types"
	"github.com/MarkSaravi/drone-go/utils"
)

type ImuModule struct {
	dev                         types.ImuDevice
	accData                     types.SensorData
	rotations                   types.Rotations
	gyro                        types.Rotations
	startTime                   time.Time
	readTime                    time.Time
	readingInterval             time.Duration
	accLowPassFilterCoefficient float64
	lowPassFilterCoefficient    float64
	readingData                 types.ImuReadingQualities
}

func NewIMU(imuMems types.ImuDevice, config types.FlightConfig) ImuModule {
	readingInterval := time.Duration(int64(time.Second) / int64(config.ImuDataPerSecond))
	badIntervalThereshold := readingInterval + readingInterval/20
	fmt.Println(readingInterval, badIntervalThereshold)
	return ImuModule{
		dev:                         imuMems,
		readTime:                    time.Time{},
		readingInterval:             readingInterval,
		accLowPassFilterCoefficient: config.AccLowPassFilterCoefficient,
		lowPassFilterCoefficient:    config.LowPassFilterCoefficient,
		readingData: types.ImuReadingQualities{
			Total:                 0,
			BadInterval:           0,
			MaxBadInterval:        0,
			BadData:               0,
			BadIntervalThereshold: badIntervalThereshold,
		},
		accData:   types.SensorData{Data: types.XYZ{X: 0, Y: 0, Z: 1}, Error: nil},
		rotations: types.Rotations{Roll: 0, Pitch: 0, Yaw: 0},
		gyro:      types.Rotations{Roll: 0, Pitch: 0, Yaw: 0},
	}
}

func (imu ImuModule) Close() {
	imu.dev.Close()
}

func (imu *ImuModule) ResetReadingTimes() {
	imu.startTime = time.Now()
	imu.readTime = imu.startTime
	imu.rotations = types.Rotations{Roll: 0, Pitch: 0, Yaw: 0}
	imu.gyro = types.Rotations{Roll: 0, Pitch: 0, Yaw: 0}
}

func (imu *ImuModule) GetRotations() (bool, types.ImuRotations, error) {
	now := time.Now()
	diff := now.Sub(imu.readTime)
	if diff < imu.readingInterval {
		return false, types.ImuRotations{}, nil
	}
	imu.readTime = now
	accData, gyroData, _, devErr := imu.dev.ReadSensors()
	imu.accData.Data = types.XYZ{
		X: LowPassFilter(imu.accData.Data.X, accData.Data.X, imu.accLowPassFilterCoefficient),
		Y: LowPassFilter(imu.accData.Data.Y, accData.Data.Y, imu.accLowPassFilterCoefficient),
		Z: LowPassFilter(imu.accData.Data.Z, accData.Data.Z, imu.accLowPassFilterCoefficient),
	}
	accRotations := AccelerometerRotations(imu.accData.Data)
	dg := GyroChanges(gyroData.Data, diff.Nanoseconds())
	imu.gyro = GyroRotations(dg, imu.gyro)
	rotations := GyroRotations(dg, imu.rotations)
	imu.rotations = types.Rotations{
		Roll:  LowPassFilter(rotations.Roll, accRotations.Roll, imu.lowPassFilterCoefficient),
		Pitch: LowPassFilter(rotations.Pitch, accRotations.Pitch, imu.lowPassFilterCoefficient),
		Yaw:   imu.gyro.Yaw,
	}
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
	pitch := utils.RadToDeg(math.Atan2(-data.X, math.Sqrt(data.Z*data.Z+data.Y*data.Y)))
	return types.Rotations{
		Roll:  roll,
		Pitch: pitch,
		Yaw:   0,
	}
}

func LowPassFilter(prevValue float64, newValue float64, coefficient float64) float64 {
	v1 := (1 - coefficient) * prevValue
	v2 := coefficient * newValue
	// fmt.Println(v1, v2, lpfc)
	return v1 + v2
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

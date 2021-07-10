package imu

import (
	"fmt"
	"math"
	"time"

	"github.com/MarkSaravi/drone-go/types"
	"github.com/MarkSaravi/drone-go/utils"
)

type imuModule struct {
	dev                         types.ImuDevice
	accData                     types.SensorData
	rotations                   types.Rotations
	gyro                        types.Rotations
	startTime                   time.Time
	readTime                    time.Time
	readingInterval             time.Duration
	accLowPassFilterCoefficient float64
	lowPassFilterCoefficient    float64
}

func NewIMU(imuMems types.ImuDevice, config types.ImuConfig) imuModule {
	readingInterval := time.Duration(int64(time.Second) / int64(config.ImuDataPerSecond))
	fmt.Println("reading interval: ", readingInterval)
	return imuModule{
		dev:                         imuMems,
		readTime:                    time.Time{},
		readingInterval:             readingInterval,
		accLowPassFilterCoefficient: config.AccLowPassFilterCoefficient,
		lowPassFilterCoefficient:    config.LowPassFilterCoefficient,
		accData:                     types.SensorData{Data: types.XYZ{X: 0, Y: 0, Z: 1}, Error: nil},
		rotations:                   types.Rotations{Roll: 0, Pitch: 0, Yaw: 0},
		gyro:                        types.Rotations{Roll: 0, Pitch: 0, Yaw: 0},
	}
}

func (imu imuModule) Close() {
	imu.dev.Close()
}

func (imu *imuModule) ResetReadingTimes() {
	imu.startTime = time.Now()
	imu.readTime = imu.startTime
	imu.rotations = types.Rotations{Roll: 0, Pitch: 0, Yaw: 0}
	imu.gyro = types.Rotations{Roll: 0, Pitch: 0, Yaw: 0}
}

func (imu *imuModule) CanRead() bool {
	if time.Since(imu.readTime) >= imu.readingInterval {
		return true
	}
	return false
}

func (imu *imuModule) GetRotations() (types.ImuRotations, error) {
	now := time.Now()
	diff := now.Sub(imu.readTime)
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

	return types.ImuRotations{
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

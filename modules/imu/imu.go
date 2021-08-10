package imu

import (
	"fmt"
	"math"
	"time"
)

// XYZ is X, Y, Z data
type XYZ struct {
	X, Y, Z float64
}

type SensorData struct {
	Error error
	Data  XYZ
}

type Rotations struct {
	Roll, Pitch, Yaw float64
}

type ImuRotations struct {
	Accelerometer Rotations
	Gyroscope     Rotations
	Rotations     Rotations
	ReadTime      int64
	ReadInterval  int64
}

type RotationsChanges struct {
	DRoll, DPitch, DYaw float64
}

type ImuMems interface {
	ReadSensors() (
		SensorData,
		SensorData,
		SensorData,
		error)
}

// IMU is interface to imu mems
type IMU interface {
	GetRotations() (ImuRotations, error)
	ResetReadingTimes()
	CanRead() bool
}

type ImuConfig struct {
	ImuDataPerSecond            int     `yaml:"imu_data_per_second"`
	AccLowPassFilterCoefficient float64 `yaml:"acc_lowpass_filter_coefficient"`
	LowPassFilterCoefficient    float64 `yaml:"lowpass_filter_coefficient"`
}

type imuModule struct {
	dev                         ImuMems
	accData                     SensorData
	rotations                   Rotations
	gyro                        Rotations
	startTime                   time.Time
	readTime                    time.Time
	readingInterval             time.Duration
	accLowPassFilterCoefficient float64
	lowPassFilterCoefficient    float64
}

func CreateIM(imuMems ImuMems, config ImuConfig) IMU {
	readingInterval := time.Duration(int64(time.Second) / int64(config.ImuDataPerSecond))
	fmt.Println("reading interval: ", readingInterval)
	imudevice := imuModule{
		dev:                         imuMems,
		readTime:                    time.Time{},
		readingInterval:             readingInterval,
		accLowPassFilterCoefficient: config.AccLowPassFilterCoefficient,
		lowPassFilterCoefficient:    config.LowPassFilterCoefficient,
		accData:                     SensorData{Data: XYZ{X: 0, Y: 0, Z: 1}, Error: nil},
		rotations:                   Rotations{Roll: 0, Pitch: 0, Yaw: 0},
		gyro:                        Rotations{Roll: 0, Pitch: 0, Yaw: 0},
	}
	return &imudevice
}

func (imu *imuModule) ResetReadingTimes() {
	imu.startTime = time.Now()
	imu.readTime = imu.startTime
	imu.rotations = Rotations{Roll: 0, Pitch: 0, Yaw: 0}
	imu.gyro = Rotations{Roll: 0, Pitch: 0, Yaw: 0}
}

func (imu *imuModule) CanRead() bool {
	if time.Since(imu.readTime) >= imu.readingInterval {
		return true
	}
	return false
}

func (imu *imuModule) GetRotations() (ImuRotations, error) {
	now := time.Now()
	diff := now.Sub(imu.readTime)
	imu.readTime = now
	accData, gyroData, _, devErr := imu.dev.ReadSensors()
	imu.accData.Data = XYZ{
		X: lowPassFilter(imu.accData.Data.X, accData.Data.X, imu.accLowPassFilterCoefficient),
		Y: lowPassFilter(imu.accData.Data.Y, accData.Data.Y, imu.accLowPassFilterCoefficient),
		Z: lowPassFilter(imu.accData.Data.Z, accData.Data.Z, imu.accLowPassFilterCoefficient),
	}
	accRotations := calcaAcelerometerRotations(imu.accData.Data)
	dg := gyroChanges(gyroData.Data, diff.Nanoseconds())
	imu.gyro = calcGyroRotations(dg, imu.gyro)
	rotations := calcGyroRotations(dg, imu.rotations)
	imu.rotations = Rotations{
		Roll:  lowPassFilter(rotations.Roll, accRotations.Roll, imu.lowPassFilterCoefficient),
		Pitch: lowPassFilter(rotations.Pitch, accRotations.Pitch, imu.lowPassFilterCoefficient),
		Yaw:   imu.gyro.Yaw,
	}

	return ImuRotations{
		Accelerometer: accRotations,
		Gyroscope:     imu.gyro,
		Rotations:     imu.rotations,
		ReadTime:      imu.readTime.UnixNano() - imu.startTime.UnixNano(),
		ReadInterval:  diff.Nanoseconds(),
	}, devErr
}

func gyroChanges(gyro XYZ, timeInterval int64) RotationsChanges {
	dt := goDurToDt(timeInterval)
	return RotationsChanges{
		DRoll:  gyro.X * dt,
		DPitch: gyro.Y * dt,
		DYaw:   gyro.Z * dt,
	}
}

func calcGyroRotations(dGyro RotationsChanges, gyro Rotations) Rotations {
	return Rotations{
		Roll:  math.Mod(gyro.Roll+dGyro.DRoll, 360),
		Pitch: math.Mod(gyro.Pitch+dGyro.DPitch, 360),
		Yaw:   math.Mod(gyro.Yaw+dGyro.DYaw, 360),
	}
}

func calcaAcelerometerRotations(data XYZ) Rotations {
	roll := radToDeg(math.Atan2(data.Y, data.Z))
	pitch := radToDeg(math.Atan2(-data.X, math.Sqrt(data.Z*data.Z+data.Y*data.Y)))
	return Rotations{
		Roll:  roll,
		Pitch: pitch,
		Yaw:   0,
	}
}

func radToDeg(x float64) float64 {
	return x / math.Pi * 180
}

func lowPassFilter(prevValue float64, newValue float64, coefficient float64) float64 {
	v1 := (1 - coefficient) * prevValue
	v2 := coefficient * newValue
	// fmt.Println(v1, v2, lpfc)
	return v1 + v2
}

func goDurToDt(d int64) float64 {
	return float64(d) / 1e9
}

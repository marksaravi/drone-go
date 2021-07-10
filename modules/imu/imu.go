package imu

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/MarkSaravi/drone-go/hardware/icm20948"
	"github.com/MarkSaravi/drone-go/types"
)

// imuHardware is interface for the imu mems
type imuHardware interface {
	Close()
	InitDevice() error
	ReadSensorsRawData() ([]byte, error)
	ReadSensors() (acc types.SensorData, gyro types.SensorData, mag types.SensorData, err error)
	WhoAmI() (string, byte, error)
}

type imuModule struct {
	dev                         imuHardware
	accData                     types.SensorData
	rotations                   types.Rotations
	gyro                        types.Rotations
	startTime                   time.Time
	readTime                    time.Time
	readingInterval             time.Duration
	accLowPassFilterCoefficient float64
	lowPassFilterCoefficient    float64
}

func CreateIM(config types.ApplicationConfig) types.IMU {
	dev, err := icm20948.NewICM20948Driver(config.Hardware.ICM20948)
	if err != nil {
		os.Exit(1)
	}
	dev.InitDevice()
	if err != nil {
		os.Exit(1)
	}
	imudevice := newIMU(dev, config.Flight.Imu)
	return &imudevice
}

func newIMU(imuMems imuHardware, config types.ImuConfig) imuModule {
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
		X: lowPassFilter(imu.accData.Data.X, accData.Data.X, imu.accLowPassFilterCoefficient),
		Y: lowPassFilter(imu.accData.Data.Y, accData.Data.Y, imu.accLowPassFilterCoefficient),
		Z: lowPassFilter(imu.accData.Data.Z, accData.Data.Z, imu.accLowPassFilterCoefficient),
	}
	accRotations := calcaAcelerometerRotations(imu.accData.Data)
	dg := gyroChanges(gyroData.Data, diff.Nanoseconds())
	imu.gyro = calcGyroRotations(dg, imu.gyro)
	rotations := calcGyroRotations(dg, imu.rotations)
	imu.rotations = types.Rotations{
		Roll:  lowPassFilter(rotations.Roll, accRotations.Roll, imu.lowPassFilterCoefficient),
		Pitch: lowPassFilter(rotations.Pitch, accRotations.Pitch, imu.lowPassFilterCoefficient),
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

func gyroChanges(gyro types.XYZ, timeInterval int64) types.RotationsChanges {
	dt := goDurToDt(timeInterval)
	return types.RotationsChanges{
		DRoll:  gyro.X * dt,
		DPitch: gyro.Y * dt,
		DYaw:   gyro.Z * dt,
	}
}

func calcGyroRotations(dGyro types.RotationsChanges, gyro types.Rotations) types.Rotations {
	return types.Rotations{
		Roll:  math.Mod(gyro.Roll+dGyro.DRoll, 360),
		Pitch: math.Mod(gyro.Pitch+dGyro.DPitch, 360),
		Yaw:   math.Mod(gyro.Yaw+dGyro.DYaw, 360),
	}
}

func calcaAcelerometerRotations(data types.XYZ) types.Rotations {
	roll := radToDeg(math.Atan2(data.Y, data.Z))
	pitch := radToDeg(math.Atan2(-data.X, math.Sqrt(data.Z*data.Z+data.Y*data.Y)))
	return types.Rotations{
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

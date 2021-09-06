package devices

import (
	"math"
	"time"

	"github.com/MarkSaravi/drone-go/models"
)

type rotationsChanges struct {
	dRoll, dPitch, dYaw float64
}

type imuMems interface {
	Read() (acc models.XYZ, gyro models.XYZ)
}

type imudevice struct {
	imuMems                     imuMems
	accRawData                  models.XYZ
	gyro                        models.Rotations
	rotations                   models.Rotations
	lastReading                 time.Time
	accLowPassFilterCoefficient float64
	lowPassFilterCoefficient    float64
}

func NewIMU(
	imuMems imuMems,
	readingInterval time.Duration,
	accLowPassFilterCoefficient float64,
	lowPassFilterCoefficient float64,
) *imudevice {
	return &imudevice{
		imuMems:                     imuMems,
		lastReading:                 time.Now(),
		accLowPassFilterCoefficient: accLowPassFilterCoefficient,
		lowPassFilterCoefficient:    lowPassFilterCoefficient,
	}
}

func (imu *imudevice) ResetTime() {
	imu.lastReading = time.Now()
}

func (imu *imudevice) ReadRotations() models.ImuRotations {
	now := time.Now()
	acc, gyro := imu.imuMems.Read()
	diff := now.Sub(imu.lastReading)
	imu.accRawData = models.XYZ{
		X: lowPassFilter(imu.accRawData.X, acc.X, imu.accLowPassFilterCoefficient),
		Y: lowPassFilter(imu.accRawData.Y, acc.Y, imu.accLowPassFilterCoefficient),
		Z: lowPassFilter(imu.accRawData.Z, acc.Z, imu.accLowPassFilterCoefficient),
	}
	accRotations := calcaAcelerometerRotations(imu.accRawData)
	dg := gyroChanges(gyro, diff.Nanoseconds())
	imu.gyro = calcGyroRotations(dg, imu.gyro)
	nrotations := calcGyroRotations(dg, imu.rotations)
	imu.rotations = models.Rotations{
		Roll:  lowPassFilter(nrotations.Roll, accRotations.Roll, imu.lowPassFilterCoefficient),
		Pitch: lowPassFilter(nrotations.Pitch, accRotations.Pitch, imu.lowPassFilterCoefficient),
		Yaw:   imu.gyro.Yaw,
	}
	imu.lastReading = now
	return models.ImuRotations{
		Accelerometer: accRotations,
		Gyroscope:     imu.gyro,
		Rotations:     imu.rotations,
		ReadTime:      now.UnixNano(),
		ReadInterval:  diff.Nanoseconds(),
	}
}

func calcaAcelerometerRotations(data models.XYZ) models.Rotations {
	roll := radToDeg(math.Atan2(data.Y, data.Z))
	pitch := radToDeg(math.Atan2(-data.X, math.Sqrt(data.Z*data.Z+data.Y*data.Y)))
	return models.Rotations{
		Roll:  roll,
		Pitch: pitch,
		Yaw:   0,
	}
}

func calcGyroRotations(dGyro rotationsChanges, gyro models.Rotations) models.Rotations {
	return models.Rotations{
		Roll:  math.Mod(gyro.Roll+dGyro.dRoll, 360),
		Pitch: math.Mod(gyro.Pitch+dGyro.dPitch, 360),
		Yaw:   math.Mod(gyro.Yaw+dGyro.dYaw, 360),
	}
}

func gyroChanges(gyro models.XYZ, timeInterval int64) rotationsChanges {
	dt := goDurToDt(timeInterval)
	return rotationsChanges{
		dRoll:  gyro.X * dt,
		dPitch: gyro.Y * dt,
		dYaw:   gyro.Z * dt,
	}
}

func lowPassFilter(prevValue float64, newValue float64, coefficient float64) float64 {
	v1 := (1 - coefficient) * prevValue
	v2 := coefficient * newValue
	return v1 + v2
}

func radToDeg(x float64) float64 {
	return x / math.Pi * 180
}

func goDurToDt(d int64) float64 {
	return float64(d) / 1e9
}

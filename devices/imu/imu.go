package imu

import (
	"math"
	"time"

	"github.com/marksaravi/drone-go/config"
	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/icm20948"
	"github.com/marksaravi/drone-go/models"
)

type rotationsChanges struct {
	dRoll, dPitch, dYaw float64
}

type imuMems interface {
	Read() (acc models.XYZ, gyro models.XYZ)
}

type imudevice struct {
	imuMems                        imuMems
	gyro                           models.Rotations
	rotations                      models.Rotations
	lastReading                    time.Time
	readingInterval                time.Duration
	complimentaryFilterCoefficient float64
}

func NewImuMems(
	imuMems imuMems,
	dataPerSecond int,
	complimentaryFilterCoefficient float64,
) *imudevice {
	return &imudevice{
		imuMems:                        imuMems,
		readingInterval:                time.Second / time.Duration(dataPerSecond),
		lastReading:                    time.Now(),
		complimentaryFilterCoefficient: complimentaryFilterCoefficient,
	}
}

func (imu *imudevice) ResetTime() {
	imu.lastReading = time.Now()
}

func (imu *imudevice) ReadRotations() (models.ImuRotations, bool) {
	if time.Since(imu.lastReading) < imu.readingInterval {
		return models.ImuRotations{}, false
	}
	now := time.Now()
	acc, gyro := imu.imuMems.Read()
	diff := now.Sub(imu.lastReading)
	accRotations := calcaAcelerometerRotations(acc)
	dg := gyroChanges(gyro, diff.Nanoseconds())
	imu.gyro = calcGyroRotations(dg, imu.gyro)
	gyroRotations := calcGyroRotations(dg, imu.rotations)
	imu.rotations = models.Rotations{
		Roll:  complimentaryFilter(gyroRotations.Roll, accRotations.Roll, imu.complimentaryFilterCoefficient),
		Pitch: complimentaryFilter(gyroRotations.Pitch, accRotations.Pitch, imu.complimentaryFilterCoefficient),
		Yaw:   imu.gyro.Yaw,
	}
	imu.lastReading = now
	return models.ImuRotations{
		Accelerometer: accRotations,
		Gyroscope:     imu.gyro,
		Rotations:     imu.rotations,
		ReadTime:      now,
		ReadInterval:  diff,
	}, true
}

func calcaAcelerometerRotations(data models.XYZ) models.Rotations {
	pitch := 180 * math.Atan2(data.X, math.Sqrt(data.Y*data.Y+data.Z*data.Z)) / math.Pi
	roll := 180 * math.Atan2(data.Y, math.Sqrt(data.X*data.X+data.Z*data.Z)) / math.Pi
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

func complimentaryFilter(gyroValue float64, accelerometerValue float64, coefficient float64) float64 {
	v1 := (1 - coefficient) * gyroValue
	v2 := coefficient * accelerometerValue
	return v1 + v2
}

func goDurToDt(d int64) float64 {
	return float64(d) / 1e9
}

func NewImu() *imudevice {
	configs := config.ReadConfigs().FlightControl
	imuConfig := configs.Imu
	imuSPIConn := hardware.NewSPIConnection(
		imuConfig.SPI.BusNumber,
		imuConfig.SPI.ChipSelect,
	)
	accConfig := imuConfig.Accelerometer
	gyroConfig := imuConfig.Gyroscope
	imuMems := icm20948.NewICM20948Driver(
		imuSPIConn,
		accConfig.SensitivityLevel,
		accConfig.Averaging,
		accConfig.LowPassFilterEnabled,
		accConfig.LowPassFilterConfig,
		accConfig.Offsets.X,
		accConfig.Offsets.Y,
		accConfig.Offsets.Z,
		gyroConfig.SensitivityLevel,
		gyroConfig.Averaging,
		gyroConfig.LowPassFilterEnabled,
		gyroConfig.LowPassFilterConfig,
		gyroConfig.Offsets.X,
		gyroConfig.Offsets.Y,
		gyroConfig.Offsets.Z,
		gyroConfig.Directions.X,
		gyroConfig.Directions.Y,
		gyroConfig.Directions.Z,
	)
	return NewImuMems(
		imuMems,
		imuConfig.DataPerSecond,
		imuConfig.ComplimentaryFilterCoefficient,
	)
}

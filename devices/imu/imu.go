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
		gyro:                           models.Rotations{Roll: 0, Pitch: 0, Yaw: 0},
		rotations:                      models.Rotations{Roll: 0, Pitch: 0, Yaw: 0},
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
	newRotations := calcGyroRotations(dg, imu.rotations)
	imu.rotations = models.Rotations{
		Roll:  complimentaryFilter(newRotations.Roll, accRotations.Roll, imu.complimentaryFilterCoefficient),
		Pitch: complimentaryFilter(newRotations.Pitch, accRotations.Pitch, imu.complimentaryFilterCoefficient),
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

func calcGyroRotations(dGyro rotationsChanges, prevRotations models.Rotations) models.Rotations {
	return models.Rotations{
		Roll:  math.Mod(prevRotations.Roll+dGyro.dRoll, 360),
		Pitch: math.Mod(prevRotations.Pitch+dGyro.dPitch, 360),
		Yaw:   math.Mod(prevRotations.Yaw+dGyro.dYaw, 360),
	}
}

func gyroChanges(gyro models.XYZ, timeInterval int64) rotationsChanges {
	dt := float64(timeInterval) / 1e9
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

func NewImu(configs config.FlightControlConfigs) *imudevice {
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

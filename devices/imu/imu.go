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
	Read() (models.XYZ, models.XYZ, error)
}

type imudevice struct {
	imuMems                        imuMems
	gyro                           models.RotationsAroundAxis
	rotations                      models.RotationsAroundAxis
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
		gyro:                           models.RotationsAroundAxis{X: 0, Y: 0, Z: 0},
		rotations:                      models.RotationsAroundAxis{X: 0, Y: 0, Z: 0},
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
	acc, gyro, err := imu.imuMems.Read()
	if err != nil {
		return models.ImuRotations{}, false
	}
	diff := now.Sub(imu.lastReading)
	accRotations := calcaAcelerometerRotations(acc)
	dg := gyroChanges(gyro, diff.Nanoseconds())
	imu.gyro = calcGyroRotations(dg, imu.gyro)
	newRotations := calcGyroRotations(dg, imu.rotations)
	imu.rotations = models.RotationsAroundAxis{
		X: complimentaryFilter(newRotations.X, accRotations.X, imu.complimentaryFilterCoefficient),
		Y: complimentaryFilter(newRotations.Y, accRotations.Y, imu.complimentaryFilterCoefficient),
		Z: imu.gyro.Z,
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

func calcaAcelerometerRotations(data models.XYZ) models.RotationsAroundAxis {
	yrot := 180 * math.Atan2(data.X, math.Sqrt(data.Y*data.Y+data.Z*data.Z)) / math.Pi
	xrot := 180 * math.Atan2(data.Y, math.Sqrt(data.X*data.X+data.Z*data.Z)) / math.Pi
	return models.RotationsAroundAxis{
		X: xrot,
		Y: yrot,
		Z: 0,
	}
}

func calcGyroRotations(dGyro rotationsChanges, prevRotations models.RotationsAroundAxis) models.RotationsAroundAxis {
	return models.RotationsAroundAxis{
		X: math.Mod(prevRotations.X+dGyro.dRoll, 360),
		Y: math.Mod(prevRotations.Y+dGyro.dPitch, 360),
		Z: math.Mod(prevRotations.Z+dGyro.dYaw, 360),
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

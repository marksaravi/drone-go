package utils

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/config"
	"github.com/MarkSaravi/drone-go/devices"
	"github.com/MarkSaravi/drone-go/devices/udplogger"
	"github.com/MarkSaravi/drone-go/drivers"
	"github.com/MarkSaravi/drone-go/drivers/icm20948"
	"github.com/MarkSaravi/drone-go/models"
)

func NewImu() interface {
	Read() (models.ImuRotations, bool)
} {
	appconfig := config.ReadFlightControlConfig()
	imuConfig := appconfig.Configs.Imu
	fmt.Println(appconfig)
	drivers.InitHost()
	imuSPIConn := drivers.NewSPIConnection(
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
	)
	readingInterval := time.Second / time.Duration(imuConfig.ImuDataPerSecond)
	imu := devices.NewIMU(
		imuMems,
		readingInterval,
		imuConfig.AccLowPassFilterCoefficient,
		imuConfig.LowPassFilterCoefficient,
	)
	return imu
}

func NewLogger() interface {
	Send(models.ImuRotations)
} {
	flightControlConfig := config.ReadFlightControlConfig()
	loggerConfig := config.ReadLoggerConfig()
	loggerConfigs := loggerConfig.UdpLoggerConfigs
	udplogger := udplogger.NewUdpLogger(
		loggerConfigs.Enabled,
		loggerConfigs.IP,
		loggerConfigs.Port,
		loggerConfigs.PacketsPerSecond,
		loggerConfigs.MaxDataPerPacket,
		flightControlConfig.Configs.Imu.ImuDataPerSecond,
	)
	return udplogger
}

func NewPowerBreaker() interface {
	SetLow()
	SetHigh()
} {
	flightControlConfig := config.ReadFlightControlConfig()
	powerbreaker := drivers.NewGPIOOutput(flightControlConfig.Configs.PowerBreaker)
	return powerbreaker
}

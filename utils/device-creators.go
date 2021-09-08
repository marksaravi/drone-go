package utils

import (
	"fmt"
	"log"

	"github.com/MarkSaravi/drone-go/config"
	"github.com/MarkSaravi/drone-go/devices"
	"github.com/MarkSaravi/drone-go/devices/motors"
	"github.com/MarkSaravi/drone-go/devices/udplogger"
	"github.com/MarkSaravi/drone-go/drivers"
	"github.com/MarkSaravi/drone-go/drivers/icm20948"
	"github.com/MarkSaravi/drone-go/drivers/nrf204"
	"github.com/MarkSaravi/drone-go/drivers/pca9685"
	"github.com/MarkSaravi/drone-go/models"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
)

func NewImu() (interface {
	ReadRotations() models.ImuRotations
	ResetTime()
}, int, int) {
	flightControlConfigs := config.ReadFlightControlConfig()
	imuConfig := flightControlConfigs.Configs.Imu
	fmt.Println(flightControlConfigs)
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
	imu := devices.NewIMU(
		imuMems,
		imuConfig.AccLowPassFilterCoefficient,
		imuConfig.LowPassFilterCoefficient,
	)
	return imu, flightControlConfigs.Configs.ImuDataPerSecond, flightControlConfigs.Configs.EscUpdatePerSecond
}

func NewLogger() interface {
	Send(models.ImuRotations)
} {
	flightControl := config.ReadFlightControlConfig()
	loggerConfig := config.ReadLoggerConfig()
	loggerConfigs := loggerConfig.UdpLoggerConfigs
	udplogger := udplogger.NewUdpLogger(
		loggerConfigs.Enabled,
		loggerConfigs.IP,
		loggerConfigs.Port,
		loggerConfigs.PacketsPerSecond,
		loggerConfigs.MaxDataPerPacket,
		flightControl.Configs.ImuDataPerSecond,
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

func NewESC() interface {
	On()
	Off()
	SetThrottles(map[uint8]float32)
} {
	flightControlConfigs := config.ReadFlightControlConfig()
	escConfigs := flightControlConfigs.Configs.ESC
	powerbreaker := drivers.NewGPIOOutput(flightControlConfigs.Configs.PowerBreaker)
	b, _ := i2creg.Open(escConfigs.I2CDev)
	i2cConn := &i2c.Dev{Addr: pca9685.PCA9685Address, Bus: b}
	pwmDev, err := pca9685.NewPCA9685(pca9685.PCA9685Address, i2cConn, escConfigs.MaxThrottle)
	if err != nil {
		log.Fatal(err)
	}
	esc := motors.NewMotorsControl(pwmDev, powerbreaker, escConfigs.MotorESCMappings)
	return esc
}

func NewRadio() interface {
	ReceiverOn()
	Receive() ([]byte, bool)
	TransmitterOn()
	Transmit([]byte) error
} {
	flightControlConfig := config.ReadFlightControlConfig()
	radioConfig := flightControlConfig.Configs.Radio
	radioSPIConn := drivers.NewSPIConnection(
		radioConfig.SPI.BusNumber,
		radioConfig.SPI.ChipSelect,
	)
	radio := nrf204.NewNRF204(radioConfig.RxTxAddress, radioConfig.CE, radioConfig.PowerDBm, radioSPIConn)
	return radio
}

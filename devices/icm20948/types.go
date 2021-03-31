package icm20948

import (
	"github.com/MarkSaravi/drone-go/types"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/host/sysfs"
)

// Register is the address and bank of the register
type Register struct {
	address byte
	bank    byte
}

type threeAxis struct {
	data     types.XYZ
	prevData types.XYZ
	dataDiff float64
	config   types.Config
}

// DeviceConfig is the configuration for the device
type DeviceConfig struct {
}

// AccelerometerConfig is the configurations for Accelerometer
type AccelerometerConfig struct {
	Sensitivity int
}

// GyroscopeConfig is the configuration for Gyroscope
type GyroscopeConfig struct {
	FullScale int
}

// Device is icm20948 mems
type Device struct {
	*sysfs.SPI
	spi.Conn
	regbank     byte
	lastReading int64
	duration    int64
	config      DeviceConfig
	acc         threeAxis
	gyro        threeAxis
}

package icm20948

import (
	"github.com/MarkSaravi/drone-go/types"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/host/sysfs"
)

// Register is the address and bank of the register
type Register struct {
	address uint8
	bank    uint8
}

// DeviceConfig is the configuration for the device
type DeviceConfig struct {
}

// AccelerometerConfig is the configurations for Accelerometer
type AccelerometerConfig struct {
	SensitivityLevel int
	XOffset          float64
	YOffset          float64
	ZOffset          float64
}

// GyroscopeConfig is the configuration for Gyroscope
type GyroscopeConfig struct {
	ScaleLevel             int
	LowPassFilterEnabled   bool
	LowPassFilter          int
	LowPassFilterAveraging int
	XOffset                float64
	YOffset                float64
	ZOffset                float64
}

// MagnetometerConfig is the configuration for Magnetometer
type MagnetometerConfig struct {
}

// Device is icm20948 mems
type Device struct {
	Name string
	*sysfs.SPI
	spi.Conn
	regbank     uint8
	lastReading int64
	duration    int64
	config      DeviceConfig
	acc         types.Sensor
	gyro        types.Sensor
	mag         types.Sensor
}

type Settings struct {
	BusNumber  int
	ChipSelect int
	Config     DeviceConfig
	AccConfig  AccelerometerConfig
	GyroConfig GyroscopeConfig
	MagConfig  MagnetometerConfig
}

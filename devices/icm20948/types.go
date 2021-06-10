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

// ImuDevice is icm20948 mems
type ImuDevice struct {
	Name string
	*sysfs.SPI
	spi.Conn
	regbank       uint8
	acc           types.Sensor
	gyro          types.Sensor
	mag           types.Sensor
	prevRotations types.Rotations
	prevGyro      types.Rotations
	prevReadTime  int64
	readTime      int64
}

// AccelerometerConfig is the configurations for Accelerometer
type AccelerometerConfig struct {
	SensitivityLevel int             `yaml:"sensitivity_level"`
	Offsets          []types.Offsets `yaml:"offsets"`
}

// GyroscopeConfig is the configuration for Gyroscope
type GyroscopeConfig struct {
	SensitivityLevel     int             `yaml:"sensitivity_level"`
	LowPassFilterEnabled bool            `yaml:"lowpass_filter_enabled"`
	LowPassFilterConfig  int             `yaml:"lowpass_filter_config"`
	Averaging            int             `yaml:"averaging"`
	Offsets              []types.Offsets `yaml:"offsets"`
}

// MagnetometerConfig is the configuration for Magnetometer
type MagnetometerConfig struct {
}

type Config struct {
	BusNumber  int                 `yaml:"bus_number"`
	ChipSelect int                 `yaml:"chip_select"`
	AccConfig  AccelerometerConfig `yaml:"accelerometer"`
	GyroConfig GyroscopeConfig     `yaml:"gyroscope"`
	MagConfig  MagnetometerConfig  `yaml:"magnetometer"`
}

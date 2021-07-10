package icm20948

import (
	"github.com/MarkSaravi/drone-go/types"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/host/sysfs"
)

// register is the address and bank of the register
type register struct {
	address uint8
	bank    uint8
}

// memsICM20948 is icm20948 mems
type memsICM20948 struct {
	Name string
	*sysfs.SPI
	spi.Conn
	regbank uint8
	acc     types.Sensor
	gyro    types.Sensor
	mag     types.Sensor
}

type offsets struct {
	X int16 `yaml:"X"`
	Y int16 `yaml:"Y"`
	Z int16 `yaml:"Z"`
}

// accelerometerConfig is the configurations for Accelerometer
type accelerometerConfig struct {
	SensitivityLevel     string  `yaml:"sensitivity_level"`
	LowPassFilterEnabled bool    `yaml:"lowpass_filter_enabled"`
	LowPassFilterConfig  int     `yaml:"lowpass_filter_config"`
	Averaging            int     `yaml:"averaging"`
	Offsets              offsets `yaml:"offsets"`
}

// gyroscopeConfig is the configuration for Gyroscope
type gyroscopeConfig struct {
	SensitivityLevel     string  `yaml:"sensitivity_level"`
	LowPassFilterEnabled bool    `yaml:"lowpass_filter_enabled"`
	LowPassFilterConfig  int     `yaml:"lowpass_filter_config"`
	Averaging            int     `yaml:"averaging"`
	Offsets              offsets `yaml:"offsets"`
}

// magnetometerConfig is the configuration for Magnetometer
type magnetometerConfig struct {
}

type Icm20948Config struct {
	BusNumber  int                 `yaml:"bus_number"`
	ChipSelect int                 `yaml:"chip_select"`
	AccConfig  accelerometerConfig `yaml:"accelerometer"`
	GyroConfig gyroscopeConfig     `yaml:"gyroscope"`
	MagConfig  magnetometerConfig  `yaml:"magnetometer"`
}

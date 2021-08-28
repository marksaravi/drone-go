package flightcontrolconfig

import (
	"log"

	"github.com/MarkSaravi/drone-go/config"
	"gopkg.in/yaml.v3"
)

type Offsets struct {
	X int16 `yaml:"X"`
	Y int16 `yaml:"Y"`
	Z int16 `yaml:"Z"`
}
type AccelerometerConfig struct {
	SensitivityLevel     string  `yaml:"sensitivity-level"`
	LowPassFilterEnabled bool    `yaml:"lowpass-filter-enabled"`
	LowPassFilterConfig  int     `yaml:"lowpass-filter-config"`
	Averaging            int     `yaml:"averaging"`
	Offsets              Offsets `yaml:"offsets"`
}

type GyroscopeConfig struct {
	SensitivityLevel     string  `yaml:"sensitivity-level"`
	LowPassFilterEnabled bool    `yaml:"lowpass-filter-enabled"`
	LowPassFilterConfig  int     `yaml:"lowpass-filter-config"`
	Averaging            int     `yaml:"averaging"`
	Offsets              Offsets `yaml:"offsets"`
}

type MagnetometerConfig struct {
	SensitivityLevel string `yaml:"sensitivity-level"`
}

type ImuMemes struct {
	Name          string              `yaml:"name"`
	SPI           config.SPI          `yaml:"spi"`
	Accelerometer AccelerometerConfig `yaml:"accelerometer"`
	Gyroscope     GyroscopeConfig     `yaml:"gyroscope"`
	Magnetometer  MagnetometerConfig  `yaml:"magnetometer"`
}

type Drivers struct {
	ImuMemes ImuMemes `yaml:"imu-mems"`
}

type ImuConfig struct {
	ImuDataPerSecond            int     `yaml:"imu-data-per-second"`
	AccLowPassFilterCoefficient float64 `yaml:"acc-lowpass-filter-coefficient"`
	LowPassFilterCoefficient    float64 `yaml:"lowpass-filter-coefficient"`
}

type Devices struct {
	ImuConfig ImuConfig `yaml:"imu"`
}

type Configs struct {
	Drivers Drivers `yaml:"drivers"`
	Devices Devices `yaml:"devices"`
}

type configs struct {
	Configs Configs `yaml:"flight-control"`
}

func ReadFlightControlConfig() configs {
	var flightcontrolconfig configs
	content := config.ReadYamlConfig()
	err := yaml.Unmarshal([]byte(content), &flightcontrolconfig)
	if err != nil {
		log.Fatalf("cannot unmarshal config: %v", err)
	}
	return flightcontrolconfig
}

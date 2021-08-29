package config

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
	SPI                         SPI                 `yaml:"spi"`
	Accelerometer               AccelerometerConfig `yaml:"accelerometer"`
	Gyroscope                   GyroscopeConfig     `yaml:"gyroscope"`
	Magnetometer                MagnetometerConfig  `yaml:"magnetometer"`
	ImuDataPerSecond            int                 `yaml:"data-per-second"`
	AccLowPassFilterCoefficient float64             `yaml:"acc-lowpass-filter-coefficient"`
	LowPassFilterCoefficient    float64             `yaml:"lowpass-filter-coefficient"`
}

type ImuConfig struct {
}

type FlightControlConfigs struct {
	Imu ImuMemes `yaml:"imu"`
}

type flightControlConfigs struct {
	Configs FlightControlConfigs `yaml:"flight-control"`
}

func ReadFlightControlConfig() flightControlConfigs {
	return readConfig(flightControlConfigs{}).(flightControlConfigs)
}

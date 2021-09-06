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
	AccLowPassFilterCoefficient float64             `yaml:"acc-lowpass-filter-coefficient"`
	LowPassFilterCoefficient    float64             `yaml:"lowpass-filter-coefficient"`
}

type ESC struct {
	I2CDev           string      `yaml:"i2c-dev"`
	MaxThrottle      float32     `yaml:"max-throttle"`
	MotorESCMappings map[int]int `yaml:"motors-esc-mappings"`
}

type Radio struct {
	CE          string `yaml:"ce-gpio"`
	RxTxAddress string `yaml:"rx-tx-address"`
	PowerDBm    string `yaml:"power-dbm"`
	SPI         SPI    `yaml:"spi"`
}

type FlightControlConfigs struct {
	ImuDataPerSecond int      `yaml:"imu-data-per-second"`
	Imu              ImuMemes `yaml:"imu"`
	ESC              ESC      `yaml:"esc"`
	Radio            Radio    `yaml:"radio"`
	PowerBreaker     string   `yaml:"power-breaker"`
}

type flightControlConfigs struct {
	Configs FlightControlConfigs `yaml:"flight-control"`
}

func ReadFlightControlConfig() flightControlConfigs {
	return readConfig(flightControlConfigs{}).(flightControlConfigs)
}

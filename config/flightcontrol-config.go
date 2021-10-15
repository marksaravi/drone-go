package config

type offsets struct {
	X int16 `yaml:"X"`
	Y int16 `yaml:"Y"`
	Z int16 `yaml:"Z"`
}

type flightControlConfigs struct {
	Configs struct {
		ImuDataPerSecond   int `yaml:"imu-data-per-second"`
		EscUpdatePerSecond int `yaml:"esc-update-per-second"`
		PID                struct {
			PGain                 float64 `yaml:"p-gain"`
			IGain                 float64 `yaml:"i-gain"`
			DGain                 float64 `yaml:"d-gain"`
			AnalogInputToThrottle float64 `yaml:"analog-input-to-throttle"`
		} `yaml:"pid"`
		Imu struct {
			SPI           spiConfig `yaml:"spi"`
			Accelerometer struct {
				SensitivityLevel     string  `yaml:"sensitivity-level"`
				LowPassFilterEnabled bool    `yaml:"lowpass-filter-enabled"`
				LowPassFilterConfig  int     `yaml:"lowpass-filter-config"`
				Averaging            int     `yaml:"averaging"`
				Offsets              offsets `yaml:"offsets"`
			} `yaml:"accelerometer"`
			Gyroscope struct {
				SensitivityLevel     string  `yaml:"sensitivity-level"`
				LowPassFilterEnabled bool    `yaml:"lowpass-filter-enabled"`
				LowPassFilterConfig  int     `yaml:"lowpass-filter-config"`
				Averaging            int     `yaml:"averaging"`
				Offsets              offsets `yaml:"offsets"`
			} `yaml:"gyroscope"`
			Magnetometer struct {
				SensitivityLevel string `yaml:"sensitivity-level"`
			} `yaml:"magnetometer"`
			AccLowPassFilterCoefficient float64 `yaml:"acc-lowpass-filter-coefficient"`
			LowPassFilterCoefficient    float64 `yaml:"lowpass-filter-coefficient"`
		} `yaml:"imu"`
		ESC struct {
			I2CDev           string      `yaml:"i2c-dev"`
			MaxThrottle      float32     `yaml:"max-throttle"`
			MotorESCMappings map[int]int `yaml:"motors-esc-mappings"`
		} `yaml:"esc"`
		Radio struct {
			CE          string    `yaml:"ce-gpio"`
			RxTxAddress string    `yaml:"rx-tx-address"`
			PowerDBm    string    `yaml:"power-dbm"`
			SPI         spiConfig `yaml:"spi"`
		} `yaml:"radio"`
		PowerBreaker string `yaml:"power-breaker"`
	} `yaml:"flight-control"`
	Radio radioConnection `yaml:"radio-connetion"`
}

func ReadFlightControlConfig() flightControlConfigs {
	return readConfig(flightControlConfigs{}).(flightControlConfigs)
}

package config

type offsets struct {
	X int16 `yaml:"X"`
	Y int16 `yaml:"Y"`
	Z int16 `yaml:"Z"`
}

type analogToDigitalConversion struct {
	Ratio  float64 `yaml:"ratio"`
	Offset float64 `yaml:"offset"`
}

type flightControlConfigs struct {
	Configs struct {
		ImuDataPerSecond   int `yaml:"imu-data-per-second"`
		EscUpdatePerSecond int `yaml:"esc-update-per-second"`
		PID                struct {
			PGain                 float64                   `yaml:"p-gain"`
			IGain                 float64                   `yaml:"i-gain"`
			DGain                 float64                   `yaml:"d-gain"`
			AnalogInputToRoll     analogToDigitalConversion `yaml:"analog-input-to-roll-conversion"`
			AnalogInputToPitch    analogToDigitalConversion `yaml:"analog-input-to-pitch-conversion"`
			AnalogInputToYaw      analogToDigitalConversion `yaml:"analog-input-to-yaw-conversion"`
			AnalogInputToThrottle analogToDigitalConversion `yaml:"analog-input-to-throttle-conversion"`
		} `yaml:"pid"`
		Imu struct {
			SPI           SpiConfig `yaml:"spi"`
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
		Radio        RadioConfig `yaml:"radio"`
		PowerBreaker string      `yaml:"power-breaker"`
	} `yaml:"flight-control"`
	RadioConnection RadioConnection `yaml:"radio-connection"`
}

func ReadFlightControlConfig() flightControlConfigs {
	return readConfig(flightControlConfigs{}).(flightControlConfigs)
}

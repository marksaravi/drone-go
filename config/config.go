package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type spiConfig struct {
	BusNumber  int `yaml:"bus-number"`
	ChipSelect int `yaml:"chip-select"`
}

type radioConfig struct {
	RxTxAddress         string    `yaml:"rx-tx-address"`
	ConnectionTimeoutMs int       `yaml:"connection-timeout-ms"`
	CE                  string    `yaml:"ce-gpio"`
	SPI                 spiConfig `yaml:"spi"`
}

type pidConfig struct {
	PGain     float64 `yaml:"p-gain"`
	IGain     float64 `yaml:"i-gain"`
	DGain     float64 `yaml:"d-gain"`
	MaxIRatio float64 `yaml:"max-i-to-max-throttle-ratio"`
}
type FlightControlConfigs struct {
	Debug                   bool    `yaml:"debug"`
	MinPIDThrottle          float64 `yaml:"min-pid-throttle"`
	MaxThrottle             float64 `yaml:"max-throttle"`
	Arm_0_2_ThrottleEnabled bool    `yaml:"arm-0-2-throttle-enabled"`
	Arm_1_3_ThrottleEnabled bool    `yaml:"arm-1-3-throttle-enabled"`
	MaxRoll                 float64 `yaml:"max-roll"`
	MaxPitch                float64 `yaml:"max-pitch"`
	MaxYaw                  float64 `yaml:"max-yaw"`

	PID struct {
		Roll        pidConfig `yaml:"rotation-around-imu-axis-x"`
		Pitch       pidConfig `yaml:"rotation-around-imu-axis-y"`
		Yaw         pidConfig `yaml:"rotation-around-imu-axis-z"`
		Calibration struct {
			Calibrating string  `yaml:"calibrating"`
			Gain        string  `yaml:"gain"`
			PStep       float64 `yaml:"p-step"`
			IStep       float64 `yaml:"i-step"`
			DStep       float64 `yaml:"d-step"`
		} `yaml:"calibration"`
	} `yaml:"pid"`

	Imu struct {
		DataPerSecond                  int     `yaml:"data-per-second"`
		ComplimentaryFilterCoefficient float64 `yaml:"complimentary-filter-coefficient"`
		Accelerometer                  struct {
			SensitivityLevel     string  `yaml:"sensitivity-level"`
			LowPassFilterEnabled bool    `yaml:"lowpass-filter-enabled"`
			LowPassFilterConfig  int     `yaml:"lowpass-filter-config"`
			Averaging            int     `yaml:"averaging"`
			Offsets              offsets `yaml:"offsets"`
		} `yaml:"accelerometer"`

		Gyroscope struct {
			SensitivityLevel     string     `yaml:"sensitivity-level"`
			LowPassFilterEnabled bool       `yaml:"lowpass-filter-enabled"`
			LowPassFilterConfig  int        `yaml:"lowpass-filter-config"`
			Averaging            int        `yaml:"averaging"`
			Offsets              offsets    `yaml:"offsets"`
			Directions           directions `yaml:"directions"`
		} `yaml:"gyroscope"`

		Magnetometer struct {
			SensitivityLevel string `yaml:"sensitivity-level"`
		} `yaml:"magnetometer"`

		SPI spiConfig `yaml:"spi"`
	} `yaml:"imu"`

	ESC struct {
		I2CDev                 string      `yaml:"i2c-dev"`
		PwmDeviceToESCMappings map[int]int `yaml:"pwm-device-to-esc-mappings"`
		UpdatePerSecond        int         `yaml:"update-per-second"`
	} `yaml:"esc"`

	Radio            radioConfig `yaml:"radio"`
	CommandPerSecond int         `yaml:"command-per-sec"`
	PowerBreaker     string      `yaml:"power-breaker"`
}

type joystick struct {
	Channel int `yaml:"channel"`
	Offset  int `yaml:"offset"`
	Dir     int `yaml:"dir"`
}

type RemoteControlConfigs struct {
	CommandPerSecond int `yaml:"command-per-sec"`

	Joysticks struct {
		Roll     joystick  `yaml:"roll"`
		Pitch    joystick  `yaml:"pitch"`
		Yaw      joystick  `yaml:"yaw"`
		Throttle joystick  `yaml:"throttle"`
		SPI      spiConfig `yaml:"spi"`
	} `yaml:"joysticks"`

	Buttons struct {
		FrontLeft   string `yaml:"front-left"`
		FrontRight  string `yaml:"front-right"`
		TopLeft     string `yaml:"top-left"`
		TopRight    string `yaml:"top-right"`
		BottomLeft  string `yaml:"bottom-left"`
		BottomRight string `yaml:"bottom-right"`
	} `yaml:"buttons"`

	Radio          radioConfig `yaml:"radio"`
	DisplayAddress uint16      `yaml:"display-i2c-address"`
	BuzzerGPIO     string      `yaml:"buzzer-gpio"`
}

type offsets struct {
	X uint16 `yaml:"X"`
	Y uint16 `yaml:"Y"`
	Z uint16 `yaml:"Z"`
}

type directions struct {
	X float64 `yaml:"X"`
	Y float64 `yaml:"Y"`
	Z float64 `yaml:"Z"`
}

type UdpLoggerConfigs struct {
	Enabled          bool   `yaml:"enabled"`
	IP               string `yaml:"ip"`
	Port             int    `yaml:"port"`
	PacketsPerSecond int    `yaml:"packets-per-second"`
}

type Configs struct {
	RemoteControl RemoteControlConfigs `yaml:"remote-control"`
	FlightControl FlightControlConfigs `yaml:"flight-control"`
	UdpLogger     UdpLoggerConfigs     `yaml:"logger"`
}

func ReadConfigs() Configs {
	content, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	var configs Configs
	yaml.Unmarshal([]byte(content), &configs)
	return configs
}

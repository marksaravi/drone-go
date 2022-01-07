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
	RxTxAddress        string    `yaml:"rx-tx-address"`
	PowerDBm           string    `yaml:"power-dbm"`
	HeartBeatTimeoutMS int       `yaml:"heart-beat-timeout-ms"`
	CE                 string    `yaml:"ce-gpio"`
	SPI                spiConfig `yaml:"spi"`
}

type flightControl struct {
	Debug             bool    `yaml:"debug"`
	MaxThrottle       float32 `yaml:"max-throttle"`
	SafeStartThrottle float32 `yaml:"safe-start-throttle"`

	PID struct {
		RollPitchPGain     float64 `yaml:"roll-pitch-p-gain"`
		RollPitchIGain     float64 `yaml:"roll-pitch-i-gain"`
		RollPitchDGain     float64 `yaml:"roll-pitch-d-gain"`
		MaxRoll            float64 `yaml:"max-roll"`
		MaxPitch           float64 `yaml:"max-pitch"`
		YawPGain           float64 `yaml:"yaw-p-gain"`
		YawIGain           float64 `yaml:"yaw-i-gain"`
		YawDGain           float64 `yaml:"yaw-d-gain"`
		MaxYaw             float64 `yaml:"max-yaw"`
		MaxI               float64 `yaml:"max-i"`
		AxisAlignmentAngle float64 `yaml:"axis-alignment-angle"`
		CalibrationGain    string  `yaml:"calibration-gain"`
		CalibrationStep    float64 `yaml:"calibration-step"`
	} `yaml:"pid"`

	Imu struct {
		DataPerSecond               int     `yaml:"data-per-second"`
		AccLowPassFilterCoefficient float64 `yaml:"acc-lowpass-filter-coefficient"`
		LowPassFilterCoefficient    float64 `yaml:"lowpass-filter-coefficient"`
		Accelerometer               struct {
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

		SPI spiConfig `yaml:"spi"`
	} `yaml:"imu"`

	ESC struct {
		I2CDev                 string      `yaml:"i2c-dev"`
		PwmDeviceToESCMappings map[int]int `yaml:"pwm-device-to-esc-mappings"`
		UpdatePerSecond        int         `yaml:"update-per-second"`
	} `yaml:"esc"`

	Radio radioConfig `yaml:"radio"`

	PowerBreaker string `yaml:"power-breaker"`
}

type joystick struct {
	Channel int `yaml:"channel"`
	Offset  int `yaml:"offset"`
}

type remoteControl struct {
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
	X int16 `yaml:"X"`
	Y int16 `yaml:"Y"`
	Z int16 `yaml:"Z"`
}

type analogToDigitalConversion struct {
	Ratio  float64 `yaml:"ratio"`
	Offset float64 `yaml:"offset"`
}

type udpLogger struct {
	Enabled          bool   `yaml:"enabled"`
	IP               string `yaml:"ip"`
	Port             int    `yaml:"port"`
	PacketsPerSecond int    `yaml:"packets-per-second"`
}

type configs struct {
	RemoteControl remoteControl `yaml:"remote-control"`
	FlightControl flightControl `yaml:"flight-control"`
	UdpLogger     udpLogger     `yaml:"logger"`
}

func ReadConfigs() configs {
	content, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	var configs configs
	yaml.Unmarshal([]byte(content), &configs)
	return configs
}

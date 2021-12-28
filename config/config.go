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
	Radio        radioConfig `yaml:"radio"`
	PowerBreaker string      `yaml:"power-breaker"`
}

type joystick struct {
	Channel       int    `yaml:"channel"`
	DigitalOffset uint16 `yaml:"digital-offset"`
}

type remoteControl struct {
	CommandPerSecond int `yaml:"command-per-sec"`
	Joysticks        struct {
		Roll       joystick  `yaml:"roll"`
		Pitch      joystick  `yaml:"pitch"`
		Yaw        joystick  `yaml:"yaw"`
		Throttle   joystick  `yaml:"throttle"`
		ValueRange uint16    `yaml:"value-range"`
		SPI        spiConfig `yaml:"spi"`
	} `yaml:"joysticks"`
	Buttons struct {
		FrontLeft   string `yaml:"front-left"`
		FrontRight  string `yaml:"front-right"`
		TopLeft     string `yaml:"top-left"`
		TopRight    string `yaml:"top-right"`
		BottomLeft  string `yaml:"bottom-left"`
		BottomRight string `yaml:"bottom-right"`
	} `yaml:"buttons"`
	Radio      radioConfig `yaml:"radio"`
	BuzzerGPIO string      `yaml:"buzzer-gpio"`
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

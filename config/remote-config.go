package config

import (
	"log"

	"gopkg.in/yaml.v3"
	"periph.io/x/periph/conn/spi"
)

type Throttle struct {
	Channel int
}

type Joystick struct {
	Channel   int     `yaml:"channel"`
	ZeroValue float32 `yaml:"zero-value"`
}

type SPI struct {
	BusNumber  int      `yaml:"bus-number"`
	ChipSelect int      `yaml:"chip-select"`
	Mode       spi.Mode `yaml:"mode"`
	Speed      int      `yaml:"speed-mega-hz"`
}

type Joysticks struct {
	Roll  Joystick `yaml:"roll"`
	Pitch Joystick `yaml:"pitch"`
	Yaw   Joystick `yaml:"yaw"`
	VRef  float32  `yaml:"v-ref"`
	SPI   SPI      `yaml:"spi"`
}

type Buttons struct {
	FrontLeft   string `yaml:"front-left"`
	FrontRight  string `yaml:"front-right"`
	TopLeft     string `yaml:"top-left"`
	TopRight    string `yaml:"top-right"`
	BottomLeft  string `yaml:"bottom-left"`
	BottomRight string `yaml:"bottom-right"`
}

type RemoteConfig struct {
	Joysticks Joysticks `yaml:"joysticks"`
	Buttons   Buttons   `yaml:"buttons"`
	Throttle  Throttle  `yaml:"throttle"`
}

type remoteConfig struct {
	RemoteConfig RemoteConfig `yaml:"remote-control"`
}

func ReadRemoteConfig() remoteConfig {
	var config remoteConfig
	content := readYamlConfig()
	err := yaml.Unmarshal([]byte(content), &config)
	if err != nil {
		log.Fatalf("cannot unmarshal config: %v", err)
	}
	return config
}

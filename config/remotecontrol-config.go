package config

import (
	"log"

	"gopkg.in/yaml.v3"
)

type Throttle struct {
	Channel int
}

type Joystick struct {
	Channel   int     `yaml:"channel"`
	ZeroValue float32 `yaml:"zero-value"`
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

type RemoteControlConfig struct {
	Joysticks Joysticks `yaml:"joysticks"`
	Buttons   Buttons   `yaml:"buttons"`
	Throttle  Throttle  `yaml:"throttle"`
}

type remoteControlConfig struct {
	RemoteControlConfig RemoteControlConfig `yaml:"remote-control"`
}

func ReadRemoteControlConfig() remoteControlConfig {
	var config remoteControlConfig
	content := readYamlConfig()
	err := yaml.Unmarshal([]byte(content), &config)
	if err != nil {
		log.Fatalf("cannot unmarshal config: %v", err)
	}
	return config
}

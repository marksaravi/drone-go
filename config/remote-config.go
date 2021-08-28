package config

import (
	"log"

	"gopkg.in/yaml.v3"
)

type Joystick struct {
	Channel   int     `yaml:"channel"`
	ZeroValue float32 `yaml:"zero-value"`
}

type Joysticks struct {
	Roll  Joystick `yaml:"roll"`
	Pitch Joystick `yaml:"pitch"`
	Yaw   Joystick `yaml:"yaw"`
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

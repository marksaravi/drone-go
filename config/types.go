package config

import "github.com/MarkSaravi/drone-go/connectors"

type RemoteControlConfig struct {
	SPI connectors.SPIConfig `yaml:"spi"`
}

type ApplicationConfig struct {
	RemoteControl RemoteControlConfig `yaml:"remote-control"`
}

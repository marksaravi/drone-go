package config

import "github.com/MarkSaravi/drone-go/remotecontrol"

type ApplicationConfig struct {
	RemoteControl remotecontrol.RemoteControlConfig `yaml:"remote-control"`
}

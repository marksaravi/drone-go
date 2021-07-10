package main

import (
	"github.com/MarkSaravi/drone-go/hardware"
	"github.com/MarkSaravi/drone-go/types"
)

type ApplicationConfig struct {
	Flight   types.FlightConfig      `yaml:"flight_control"`
	Hardware hardware.HardwareConfig `yaml:"devices"`
	UDP      types.UdpLoggerConfig   `yaml:"udp"`
}

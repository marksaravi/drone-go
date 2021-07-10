package main

import (
	"github.com/MarkSaravi/drone-go/hardware/icm20948"
	"github.com/MarkSaravi/drone-go/types"
)

type ApplicationConfig struct {
	Flight  types.FlightConfig `yaml:"flight_control"`
	Devices struct {
		ICM20948 icm20948.Config `yaml:"icm20948"`
	} `yaml:"devices"`
	UDP types.UdpLoggerConfig `yaml:"udp"`
}

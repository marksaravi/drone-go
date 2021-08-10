package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/MarkSaravi/drone-go/hardware/icm20948"
	"github.com/MarkSaravi/drone-go/hardware/nrf204"
	"github.com/MarkSaravi/drone-go/hardware/pca9685"
	"github.com/MarkSaravi/drone-go/modules/imu"
	"github.com/MarkSaravi/drone-go/remotecontrol"
	"gopkg.in/yaml.v3"
)

type PidConfig struct {
	ProportionalGain float32 `yaml:"proportional–gain"`
	IntegralGain     float32 `yaml:"integral-gain"`
	DerivativeGain   float32 `yaml:"derivative-gain"`
}

type EscConfig struct {
	UpdateFrequency int     `yaml:"update_frequency"`
	MaxThrottle     float32 `yaml:"max_throttle"`
}

type FlightConfig struct {
	PID PidConfig     `yaml:"pid"`
	Imu imu.ImuConfig `yaml:"imu"`
	Esc EscConfig     `yaml:"esc"`
}

type HardwareConfig struct {
	ICM20948 icm20948.ICM20948Config `yaml:"icm20948"`
	PCA9685  pca9685.PCA9685Config   `yaml:"pca9685"`
	NRF204   nrf204.NRF204Config     `yaml:"nrf204"`
}

type UdpLoggerConfig struct {
	Enabled          bool   `yaml:"enabled"`
	IP               string `yaml:"ip"`
	Port             int    `yaml:"port"`
	PacketsPerSecond int    `yaml:"packets_per_second"`
	MaxDataPerPacket int    `yaml:"max_data_per_packet"`
}

type ApplicationConfig struct {
	Flight        FlightConfig                      `yaml:"flight_control"`
	Hardware      HardwareConfig                    `yaml:"devices"`
	UDP           UdpLoggerConfig                   `yaml:"udp"`
	RemoteControl remotecontrol.RemoteControlConfig `yaml:"remote-control"`
}

func ReadConfigs() ApplicationConfig {
	var config ApplicationConfig

	content, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	err = yaml.Unmarshal([]byte(content), &config)
	if err != nil {
		log.Fatalf("cannot unmarshal config: %v", err)
		os.Exit(1)
	}
	fmt.Println(config)
	return config
}

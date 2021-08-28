package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/MarkSaravi/drone-go/apps/flightcontrol"
	"github.com/MarkSaravi/drone-go/hardware/icm20948"
	"github.com/MarkSaravi/drone-go/hardware/nrf204"
	"github.com/MarkSaravi/drone-go/hardware/pca9685"
	"github.com/MarkSaravi/drone-go/modules/udplogger"
	"gopkg.in/yaml.v3"
)

type HardwareConfig struct {
	ICM20948 icm20948.ICM20948Config `yaml:"icm20948"`
	PCA9685  pca9685.PCA9685Config   `yaml:"pca9685"`
	NRF204   nrf204.NRF204Config     `yaml:"nrf204"`
}

type ApplicationConfig struct {
	Flight   flightcontrol.FlightConfig `yaml:"flight_control"`
	Hardware HardwareConfig             `yaml:"devices"`
	UDP      udplogger.UdpLoggerConfig  `yaml:"udp"`
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

func readYamlConfig() []byte {
	content, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	return content
}

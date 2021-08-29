package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
	"periph.io/x/periph/conn/spi"
)

type SPI struct {
	BusNumber  int      `yaml:"bus-number"`
	ChipSelect int      `yaml:"chip-select"`
	Mode       spi.Mode `yaml:"mode"`
	Speed      int      `yaml:"speed-mega-hz"`
}

type UdpLogger struct {
	Enabled          bool   `yaml:"enabled"`
	IP               string `yaml:"ip"`
	Port             int    `yaml:"port"`
	PacketsPerSecond int    `yaml:"packets-per-second"`
	MaxDataPerPacket int    `yaml:"max-data-per-packet"`
}

func readConfig(config interface{}) interface{} {
	content, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	switch typedConfig := config.(type) {
	case flightControlConfigs:
		yaml.Unmarshal([]byte(content), &typedConfig)
		return typedConfig
	case remoteControlConfigs:
		yaml.Unmarshal([]byte(content), &typedConfig)
		return typedConfig
	default:
		log.Fatalf("cannot unmarshal config: undefined type")
	}
	return nil
}

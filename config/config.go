package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type SpiConfig struct {
	BusNumber  int `yaml:"bus-number"`
	ChipSelect int `yaml:"chip-select"`
}

type RadioConfig struct {
	CE  string    `yaml:"ce-gpio"`
	SPI SpiConfig `yaml:"spi"`
}

type RadioConnection struct {
	RxTxAddress         string `yaml:"rx-tx-address"`
	PowerDBm            string `yaml:"power-dbm"`
	ConnectionTimeoutMS int    `yaml:"connection-timeout-ms"`
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
	case udpLoggerConfigs:
		yaml.Unmarshal([]byte(content), &typedConfig)
		return typedConfig
	default:
		log.Fatalf("cannot unmarshal config: undefined type")
	}
	return nil
}

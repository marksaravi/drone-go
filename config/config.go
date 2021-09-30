package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type spiConfig struct {
	BusNumber  int `yaml:"bus-number"`
	ChipSelect int `yaml:"chip-select"`
}

type radioConnection struct {
	CommandPerSecond          int `yaml:"command-per-second"`
	DroneConnectionTimeoutMS  int `yaml:"drone-connection-timeout-ms"`
	RemoteConnectionTimeoutMS int `yaml:"remote-connection-timeout-ms"`
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

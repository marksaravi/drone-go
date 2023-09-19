package drone

import (
	"log"
	"os"

	"encoding/json"
)

type DroneConfigs struct {
	CommandsPerSecond int `json:"commands-per-second"`
	Radio             struct {
		RxTxAddress string `json:"rx-tx-address"`
		SPI         struct {
			BusNum             int    `json:"bus-num"`
			ChipSelect         int    `json:"chip-select"`
			SpiChipEnabledGPIO string `json:"chip-enabled-gpio"`
		}
	} `json:"radio"`
}

func ReadConfigs(configPath string) DroneConfigs {
	content, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	var configs DroneConfigs
	json.Unmarshal([]byte(content), &configs)
	return configs
}

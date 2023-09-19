package remote

import (
	"log"
	"os"

	"encoding/json"
)

type RemoteConfigs struct {
	CommandsPerSecond int `json:"commands-per-second"`
	Radio             struct {
		RxTxAddress string `json:"rx-tx-address"`
		SPI         struct {
			BusNum             int    `json:"bus-num"`
			ChipSelect         int    `json:"chip-select"`
			SpiChipEnabledGPIO string `json:"chip-enabled-gpio"`
		}
	} `json:"radio"`
	PushButtons struct {
		Right []string `json:"right"`
		Left  []string `json:"left"`
	} `json:"push-buttons-gpio"`
}

func ReadConfigs(configPath string) RemoteConfigs {
	content, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	var configs RemoteConfigs
	json.Unmarshal([]byte(content), &configs)
	return configs
}

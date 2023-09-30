package drone

import (
	"log"
	"os"

	"encoding/json"

	"github.com/marksaravi/drone-go/hardware"
)

type DroneConfigs struct {
	CommandsPerSecond int `json:"commands-per-second"`
	Radio             struct {
		RxTxAddress string                  `json:"rx-tx-address"`
		SPI         hardware.SPIConnConfigs `json:"spi"`
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

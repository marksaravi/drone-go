package remote

import (
	"log"
	"os"

	"encoding/json"

	"github.com/marksaravi/drone-go/hardware"
)

type RemoteConfigs struct {
	CommandsPerSecond int `json:"commands-per-second"`
	Radio             struct {
		RxTxAddress string                  `json:"rx-tx-address"`
		SPI         hardware.SPIConnConfigs `json:"spi"`
	} `json:"radio"`
	Joysticks struct {
		SPI             hardware.SPIConnConfigs `json:"spi"`
		RollChannel     int                     `json:"roll-channel"`
		PitchChannel    int                     `json:"pitch-channel"`
		YawChannel      int                     `json:"yaw-channel"`
		ThrottleChannel int                     `json:"throttle-channel"`
	} `json:"joysticks"`
	PushButtons []struct {
		Name         string `json:"name"`
		Index        int    `json:"index"`
		IsPushButton bool `json:"is-push-button"`
		GPIO         string `json:"gpio"`
	} `json:"push-buttons-gpio"`
	DisplayUpdatePerSecond int `json:"display-update-per-second"`
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

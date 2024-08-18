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
		RollChannel     int                     `json:"roll-channel"`
		RollMin         float64                  `json:"roll-min"`
		RollMid         float64                  `json:"roll-mid"`
		RollMax         float64                  `json:"roll-max"`

		PitchChannel    int                     `json:"pitch-channel"`
		PitchMin        float64                 `json:"pitch-min"`
		PitchMid        float64                 `json:"pitch-mid"`
		PitchMax        float64                 `json:"pitch-max"`		
		
		YawChannel      int                     `json:"yaw-channel"`
		YawMin          float64                 `json:"yaw-min"`
		YawMid          float64                 `json:"yaw-mid"`
		YawMax          float64                 `json:"yaw-max"`		

		ThrottleChannel int                     `json:"throttle-channel"`
		ThrottleMin     float64                 `json:"throttle-min"`
		ThrottleMid     float64                 `json:"throttle-mid"`
		ThrottleMax     float64                 `json:"throttle-max"`
		
		I2CAddress             hardware.SPIConnConfigs `json:"i2c-address"`
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

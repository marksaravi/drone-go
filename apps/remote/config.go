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
		RollChannel     int                    `json:"roll-channel"`
		RollMin         uint16                 `json:"roll-min"`
		RollMid         uint16                 `json:"roll-mid"`
		RollMax         uint16                 `json:"roll-max"`

		PitchChannel    int                    `json:"pitch-channel"`
		PitchMin        uint16                 `json:"pitch-min"`
		PitchMid        uint16                 `json:"pitch-mid"`
		PitchMax        uint16                 `json:"pitch-max"`		
		
		YawChannel      int                    `json:"yaw-channel"`
		YawMin          uint16                 `json:"yaw-min"`
		YawMid          uint16                 `json:"yaw-mid"`
		YawMax          uint16                 `json:"yaw-max"`		

		ThrottleChannel int                    `json:"throttle-channel"`
		ThrottleMin     uint16                 `json:"throttle-min"`
		ThrottleMid     uint16                 `json:"throttle-mid"`
		ThrottleMax     uint16                 `json:"throttle-max"`
		
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

package remote

import (
	"log"
	"os"

	"encoding/json"
)

type SPI struct {
	BusNum             int    `json:"bus-num"`
	ChipSelect         int    `json:"chip-select"`
	SpiChipEnabledGPIO string `json:"chip-enabled-gpio"`
}

type RemoteConfigs struct {
	CommandsPerSecond int `json:"commands-per-second"`
	Radio             struct {
		RxTxAddress string `json:"rx-tx-address"`
		SPI         SPI    `json:"spi"`
	}
	Joysticks struct {
		SPI             SPI    `json:"spi"`
		RollChannel     int    `json:"roll-channel"`
		PitchChannel    int    `json:"pitch-channel"`
		YawChannel      int    `json:"yaw-channel"`
		ThrottleChannel int    `json:"throttle-channel"`
		RollMidValue    uint16 `json:"roll-mid-value"`
		PitchMidValue   uint16 `json:"pitch-mid-value"`
		YawMidValue     uint16 `json:"yaw-mid-value"`
	} `json:"joysticks"`
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

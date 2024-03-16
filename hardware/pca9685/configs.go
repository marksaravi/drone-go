package pca9685

import (
	"encoding/json"
	"log"
	"os"
)

type HardwareConfigs struct {
	Configs Configs `json:"pca9685"`
}

func ReadConfigs(configPath string) Configs {
	content, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	var configs HardwareConfigs
	json.Unmarshal([]byte(content), &configs)
	return configs.Configs
}

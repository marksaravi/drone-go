package imu

import (
	"encoding/json"
	"log"
	"os"
)

type ImuConfigs struct {
	Configs Configs `json:"imu"`
}

func ReadConfigs(configPath string) Configs {
	content, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	var configs ImuConfigs
	json.Unmarshal([]byte(content), &configs)
	return configs.Configs
}

package icm20789

import (
	"encoding/json"
	"log"
	"os"
)

type HardwareConfigs struct {
	Configs Configs `json:"icm20789"`
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

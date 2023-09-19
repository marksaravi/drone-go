package remote

import (
	"fmt"
	"log"
	"os"

	"encoding/json"
)

type RemoteConfig struct {
	PushButtons struct {
		Right []string `json:"right"`
		Left  []string `json:"left"`
	} `json:"push-buttons-gpio"`
}

func ReadConfigs(configPath string) RemoteConfig {
	content, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(content))
	var configs RemoteConfig
	json.Unmarshal([]byte(content), &configs)
	return configs
}

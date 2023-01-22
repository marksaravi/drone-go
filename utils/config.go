package utils

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

func ReadConfigs() any {
	content, err := os.ReadFile("./config-2.yaml")
	if err != nil {
		log.Fatal(err)
	}
	var configs any
	yaml.Unmarshal([]byte(content), &configs)
	return configs
}

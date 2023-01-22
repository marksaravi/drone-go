package utils

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

func ReadConfigs(configs any) any {
	content, err := os.ReadFile("./config-2.yaml")
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal([]byte(content), configs)
	if err != nil {
		log.Fatal(err)
	}
	return configs
}

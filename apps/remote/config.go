package remote

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type RemoteConfig struct {
	Buttons map[string]string `yaml:"push-buttons"`
}

func ReadConfigs(configPath string) RemoteConfig {
	content, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	type Configs struct {
		Remote RemoteConfig `yaml:"remote"`
	}

	var configs Configs
	yaml.Unmarshal([]byte(content), &configs)
	return configs.Remote
}

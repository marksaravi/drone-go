package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/MarkSaravi/drone-go/types"
	"gopkg.in/yaml.v3"
)

func ReadConfigs() types.ApplicationConfig {
	var config types.ApplicationConfig

	content, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	err = yaml.Unmarshal([]byte(content), &config)
	if err != nil {
		log.Fatalf("cannot unmarshal config: %v", err)
		os.Exit(1)
	}
	fmt.Println(config)
	return config
}

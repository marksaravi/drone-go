package config

import (
	"io/ioutil"
	"log"

	"periph.io/x/periph/conn/spi"
)

type SPI struct {
	BusNumber  int      `yaml:"bus-number"`
	ChipSelect int      `yaml:"chip-select"`
	Mode       spi.Mode `yaml:"mode"`
	Speed      int      `yaml:"speed-mega-hz"`
}

func readYamlConfig() []byte {
	content, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	return content
}

package hardware

import (
	"log"

	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/host/v3/sysfs"
)

type SPIConnConfigs struct {
	BusNumber       int    `json:"bus-number"`
	ChipSelect      int    `json:"chip-select"`
	ChipEnabledGPIO string `json:"chip-enabled-gpio"`
}

func NewSPIConnection(configs SPIConnConfigs) spi.Conn {
	p, err := sysfs.NewSPI(configs.BusNumber, configs.ChipSelect)

	if err != nil {
		log.Fatal(err)
	}

	// Convert the spi.Port into a spi.Conn so it can be used for communication.
	c, err := p.Connect(physic.MegaHertz, spi.Mode0, 8)

	if err != nil {
		log.Fatal(err)
	}
	return c
}

package hardware

import (
	"log"

	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/host/v3/sysfs"
)

func NewSPIConnection(busNumber int, chipSelect int) spi.Conn {
	p, err := sysfs.NewSPI(busNumber, chipSelect)

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

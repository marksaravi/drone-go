package hardware

import (
	"log"

	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/host/sysfs"
)

func NewSPIConnection(busNumber int, chipSelect int) spi.Conn {
	spibus, _ := sysfs.NewSPI(
		busNumber,
		chipSelect,
	)
	spiConn, err := spibus.Connect(
		physic.Frequency(7)*physic.MegaHertz,
		spi.Mode0,
		8,
	)
	if err != nil {
		log.Fatal(err)
	}
	return spiConn
}

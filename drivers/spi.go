package drivers

import (
	"log"

	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/host/sysfs"
)

func NewSPIConnection(busNumber int, chipSelect int, speed int, mode spi.Mode) spi.Conn {
	spibus, _ := sysfs.NewSPI(
		busNumber,
		chipSelect,
	)
	spiConn, err := spibus.Connect(
		physic.Frequency(speed)*physic.MegaHertz,
		mode,
		8,
	)
	if err != nil {
		log.Fatal(err)
	}
	return spiConn
}

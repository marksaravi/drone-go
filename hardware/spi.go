package hardware

import (
	"log"

	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/host/sysfs"
)

type SPISettings struct {
	BusNumber  int
	ChipSelect int
}

func NewSPI(settings SPISettings) spi.Conn {
	spibus, _ := sysfs.NewSPI(
		settings.BusNumber,
		settings.ChipSelect,
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

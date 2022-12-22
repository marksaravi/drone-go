package hardware

import (
	"fmt"
	"log"
)

type physic interface{}
type spi interface{}
type sysfs interface{}

func NewSPIConnection(busNumber int, chipSelect int) spi.Conn {
	spibus, _ := sysfs.NewSPI(
		busNumber,
		chipSelect,
	)
	spiConn, err := spibus.Connect(
		physic.Frequency(1)*physic.MegaHertz,
		spi.Mode0,
		8,
	)
	if err != nil {
		log.Fatal(err)
	}
	if p, ok := spiConn.(spi.Pins); ok {
		fmt.Printf("  CLK : %s\n", p.CLK())
		fmt.Printf("  MOSI: %s\n", p.MOSI())
		fmt.Printf("  MISO: %s\n", p.MISO())
		fmt.Printf("  CS  : %s\n", p.CS())
	}

	return spiConn
}

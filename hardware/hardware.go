package hardware

import (
	"log"

	"github.com/MarkSaravi/drone-go/devices/radiolink"
	"github.com/MarkSaravi/drone-go/drivers/nrf204"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/host"
	"periph.io/x/periph/host/sysfs"
)

func NewRadioLink(config nrf204.NRF204Config) radiolink.RadioLink {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}
	spibus, err := sysfs.NewSPI(config.BusNumber, config.ChipSelect)
	if err != nil {
		log.Fatal(err)
	}
	spiconn, err := spibus.Connect(physic.MegaHertz, spi.Mode0, 8)
	if err != nil {
		log.Fatal(err)
	}
	return nrf204.NewNRF204(config, spiconn)
}

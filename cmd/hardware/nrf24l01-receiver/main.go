package main

import (
	"fmt"
	"log"
	"time"

	"github.com/MarkSaravi/drone-go/hardware/nrf204"
	"github.com/MarkSaravi/drone-go/types"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/host"
	"periph.io/x/periph/host/sysfs"
)

func main() {
	config := types.RadioLinkConfig{
		GPIO: types.RadioLinkGPIOPins{
			CE: "GPIO26",
		},
		BusNumber:  1,
		ChipSelect: 2,
		RxAddress:  "03896",
		PowerDBm:   nrf204.RF_POWER_MINUS_18dBm,
	}
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

	fmt.Println("Start")
	receiver := nrf204.CreateNRF204(config, spiconn)
	receiver.Init()
	receiver.ReceiverOn()
	var numReceive int = 0
	start := time.Now()
	for {
		if receiver.IsPayloadAvailable(0) {
			flightdata := nrf204.PayloadToFlightData(receiver.ReadPayload())
			numReceive++
			if time.Since(start) >= time.Second {
				start = time.Now()
				fmt.Println("received ", numReceive, " (", flightdata, ")")

			}
		}
	}
}

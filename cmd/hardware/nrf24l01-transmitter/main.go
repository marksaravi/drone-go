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
	receiver.StartTransmitting()
	var roll float32 = 0
	var altitude float32 = 0
	var motorsEngaged bool = false
	for range time.Tick(time.Millisecond * 1000) {
		flightdata := types.FlightData{
			Roll:          roll,
			Pitch:         -34.53,
			Yaw:           0,
			Throttle:      13.45,
			Altitude:      altitude,
			MotorsEngaged: motorsEngaged,
		}
		payload := nrf204.FlightDataToPayload(flightdata)
		fmt.Println("send (", flightdata, ")")
		err := receiver.WritePayload(payload)
		if err != nil {
			fmt.Println(err)
		}
		roll += 0.3
		altitude += 1.34
		motorsEngaged = !motorsEngaged
	}
}

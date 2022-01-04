package main

import (
	"fmt"
	"time"

	"github.com/marksaravi/drone-go/config"
	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/nrf204"
	"github.com/marksaravi/drone-go/models"
)

func main() {
	hardware.InitHost()
	radioConfigs := config.ReadConfigs().FlightControl.Radio
	radio := nrf204.NewRadio(
		radioConfigs.SPI.BusNumber,
		radioConfigs.SPI.ChipSelect,
		radioConfigs.CE,
		radioConfigs.RxTxAddress,
		radioConfigs.PowerDBm,
	)
	radio.TransmitterOn()
	var id byte = 0
	for range time.Tick(time.Second / 50) {
		var payload models.Payload
		payload[0] = id
		id++
		if id > 250 {
			id = 0
		}
		err := radio.Transmit(payload)
		if err != nil {
			fmt.Println(err)
		}
	}
}

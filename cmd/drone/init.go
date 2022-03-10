package main

import (
	"github.com/marksaravi/drone-go/config"
	"github.com/marksaravi/drone-go/devices/radio"
	"github.com/marksaravi/drone-go/hardware/nrf204"
)

func createRadioReceiver(flightcontrolConfigs config.FlightControlConfigs) *radio.RadioReceiver {
	radioConfigs := flightcontrolConfigs.Radio
	radioNRF204 := nrf204.NewNRF204EnhancedBurst(
		radioConfigs.SPI.BusNumber,
		radioConfigs.SPI.ChipSelect,
		radioConfigs.CE,
		radioConfigs.RxTxAddress,
	)
	return radio.NewReceiver(radioNRF204, flightcontrolConfigs.CommandPerSecond, radioConfigs.ConnectionTimeoutMs)
}

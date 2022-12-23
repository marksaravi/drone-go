package hardware

import (
	"log"

	"periph.io/x/host/v3"
)

func HostInitialize() {
	log.Println("SETUP HOST init")
	state, err := host.Init()
	if err != nil {
		log.Fatalf("failed to initialize periph: %v", err)
	}
	spiloaded := false
	i2cloaded := false

	for _, driver := range state.Loaded {
		if driver.String() == "sysfs-spi" {
			spiloaded = true
		}
		if driver.String() == "sysfs-i2c" {
			i2cloaded = true
		}
	}
	if !spiloaded {
		log.Fatalf("failed to initialize spi")
	}
	if !i2cloaded {
		log.Fatalf("failed to initialize i2c")
	}
}

package hardware

import (
	"log"

	"periph.io/x/conn/v3/driver/driverreg"
	"periph.io/x/host/v3"
)

func HostInitialize() {
	state, err := host.Init()
	if err != nil {
		log.Fatalf("failed to initialize periph: %v", err)
	}
	if _, err := driverreg.Init(); err != nil {
		log.Fatal(err)
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

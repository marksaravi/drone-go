package hardware

import (
	"log"

	"github.com/MarkSaravi/drone-go/hardware/nrf204"
	"github.com/MarkSaravi/drone-go/hardware/pca9685"
	"github.com/MarkSaravi/drone-go/modules/motors"
	"github.com/MarkSaravi/drone-go/modules/powerbreaker"
	"github.com/MarkSaravi/drone-go/modules/radiolink"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/host"
	"periph.io/x/periph/host/sysfs"
)

func newPowerBreaker(gpio string) powerbreaker.PowerBreaker {
	return powerbreaker.NewPowerBreaker(gpio)
}

func newPwmDev(config pca9685.PCA9685Config) motors.ESC {
	b, err := i2creg.Open(config.Device)
	d := &i2c.Dev{Addr: pca9685.PCA9685Address, Bus: b}
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
	pwmDev, err := pca9685.NewPCA9685Driver(pca9685.PCA9685Address, d, 15, config.Motors)
	if err != nil {
		log.Fatal(err)
	}
	pwmDev.Start()
	pwmDev.StopAll()
	return pwmDev
}

func newRadioLink(config nrf204.NRF204Config) radiolink.RadioLink {
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

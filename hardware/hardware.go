package hardware

import (
	"fmt"
	"log"

	"github.com/MarkSaravi/drone-go/config"
	"github.com/MarkSaravi/drone-go/connectors"
	"github.com/MarkSaravi/drone-go/connectors/i2c"
	"github.com/MarkSaravi/drone-go/hardware/icm20948"
	"github.com/MarkSaravi/drone-go/hardware/mcp3008"
	"github.com/MarkSaravi/drone-go/hardware/nrf204"
	"github.com/MarkSaravi/drone-go/hardware/pca9685"
	"github.com/MarkSaravi/drone-go/modules/adcconverter"
	"github.com/MarkSaravi/drone-go/modules/powerbreaker"
	"github.com/MarkSaravi/drone-go/types"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/host"
	"periph.io/x/periph/host/sysfs"
)

func InitDroneHardware(config types.ApplicationConfig) (types.ImuMems, types.ESC, types.RadioLink, types.PowerBreaker) {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}
	pwmDev := newPwmDev(config.Hardware.PCA9685)
	powerbreaker := newPowerBreaker(config.Hardware.PCA9685.PowerBrokerGPIO)
	imuMems := newImuMems(config.Hardware.ICM20948)
	radio := newRadioLink(config.Hardware.NRF204)
	return imuMems, pwmDev, radio, powerbreaker
}

func InitRemoteHardware(config config.ApplicationConfig) (adcconverter.AnalogToDigitalConverter, types.RadioLink) {
	fmt.Println(config)
	spibus, _ := sysfs.NewSPI(
		config.RemoteControl.MCP3008.SPI.BusNumber,
		config.RemoteControl.MCP3008.SPI.ChipSelect,
	)
	spiconn, err := spibus.Connect(
		physic.Frequency(config.RemoteControl.MCP3008.SPI.SpeedMegaHz)*physic.MegaHertz,
		connectors.ConfigToSPIMode(config.RemoteControl.MCP3008.SPI.Mode),
		8,
	)
	if err != nil {
		log.Fatal(err)
	}
	adc := mcp3008.NewMCP3008(spiconn)
	return adc, nil
}

func newImuMems(config icm20948.ICM20948Config) types.ImuMems {
	imuMems, err := icm20948.NewICM20948Driver(config)
	if err != nil {
		log.Fatal(err)
	}
	return imuMems
}

func newPowerBreaker(gpio string) types.PowerBreaker {
	return powerbreaker.NewPowerBreaker(gpio)

}

func newPwmDev(config types.PCA9685Config) types.ESC {
	i2cConnection, err := i2c.Open(config.Device)
	if err != nil {
		log.Fatal(err)
	}
	pwmDev, err := pca9685.NewPCA9685Driver(pca9685.PCA9685Address, i2cConnection, 15, config.Motors)
	if err != nil {
		log.Fatal(err)
	}
	pwmDev.Start()
	pwmDev.StopAll()
	return pwmDev
}

func newRadioLink(config types.NRF204Config) types.RadioLink {
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

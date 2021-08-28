package hardware

import (
	"log"

	"github.com/MarkSaravi/drone-go/config"
	"github.com/MarkSaravi/drone-go/hardware/icm20948"
	"github.com/MarkSaravi/drone-go/hardware/nrf204"
	"github.com/MarkSaravi/drone-go/hardware/pca9685"
	"github.com/MarkSaravi/drone-go/modules/imu"
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

func InitHost() {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}
}

func InitDroneHardware(config config.ApplicationConfig) (imu.ImuDevice, motors.ESC, radiolink.RadioLink, powerbreaker.PowerBreaker) {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}
	pwmDev := newPwmDev(config.Hardware.PCA9685)
	powerbreaker := newPowerBreaker(config.Hardware.PCA9685.PowerBrokerGPIO)
	imuDev, err := icm20948.NewICM20948Driver(config.Hardware.ICM20948)
	if err != nil {
		log.Fatal(err)
	}
	radio := newRadioLink(config.Hardware.NRF204)
	return imuDev, pwmDev, radio, powerbreaker
}

func configToSPIMode(configValue int) spi.Mode {
	switch configValue {
	case 0:
		return spi.Mode0
	case 1:
		return spi.Mode1
	case 2:
		return spi.Mode2
	case 3:
		return spi.Mode3
	default:
		return spi.Mode0
	}
}

// func InitRemoteHardware(config config.ApplicationConfig) (
// 	adcconverter.AnalogToDigitalConverter,
// 	radiolink.RadioLink,
// 	gpio.PinIn,
// 	gpio.PinIn,
// 	gpio.PinIn,
// 	gpio.PinIn,
// 	gpio.PinIn,
// 	gpio.PinIn,
// ) {
// 	if _, err := host.Init(); err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(config)
// 	spibus, _ := sysfs.NewSPI(
// 		config.RemoteControl.MCP3008.SPI.BusNumber,
// 		config.RemoteControl.MCP3008.SPI.ChipSelect,
// 	)
// 	spiconn, err := spibus.Connect(
// 		physic.Frequency(config.RemoteControl.MCP3008.SPI.SpeedMegaHz)*physic.MegaHertz,
// 		configToSPIMode(config.RemoteControl.MCP3008.SPI.Mode),
// 		8,
// 	)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	adc := mcp3008.NewMCP3008(spiconn)
// 	buttonFrontLeft := NewButton(config.RemoteControl.ButtonFrontLeft)
// 	buttonFrontRight := NewButton(config.RemoteControl.ButtonFrontRight)
// 	buttonTopLeft := NewButton(config.RemoteControl.ButtonTopLeft)
// 	buttonTopRight := NewButton(config.RemoteControl.ButtonTopRight)
// 	buttonDownLeft := NewButton(config.RemoteControl.ButtonDownLeft)
// 	buttonDownRight := NewButton(config.RemoteControl.ButtonDownRight)
// 	return adc, nil, buttonFrontLeft, buttonFrontRight, buttonTopLeft, buttonTopRight, buttonDownLeft, buttonDownRight
// }

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

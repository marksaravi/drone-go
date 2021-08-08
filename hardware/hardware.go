package hardware

import (
	"log"

	"github.com/MarkSaravi/drone-go/connectors/i2c"
	"github.com/MarkSaravi/drone-go/hardware/icm20948"
	"github.com/MarkSaravi/drone-go/hardware/pca9685"
	"github.com/MarkSaravi/drone-go/modules/powerbreaker"
	"github.com/MarkSaravi/drone-go/types"
	"periph.io/x/periph/host"
)

func InitHardware(config types.ApplicationConfig) (types.ImuMems, types.ESC, types.PowerBreaker) {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}
	i2cConnection, err := i2c.Open(config.Flight.Esc.Device)
	if err != nil {
		log.Fatal(err)
	}
	pwmDev, err := pca9685.NewPCA9685Driver(pca9685.PCA9685Address, i2cConnection, 15, config.Flight.Esc.Motors)
	if err != nil {
		log.Fatal(err)
	}
	powerbreaker := powerbreaker.NewPowerBreaker(config.Flight.Esc.PowerBrokerGPIO)
	pwmDev.Start()
	pwmDev.StopAll()
	imuMems, err := icm20948.NewICM20948Driver(config.Hardware.ICM20948)
	if err != nil {
		log.Fatal(err)
	}
	return imuMems, pwmDev, powerbreaker
}

package main

import (
	"github.com/MarkSaravi/drone-go/config"
	"github.com/MarkSaravi/drone-go/drivers"
	"github.com/MarkSaravi/drone-go/drivers/pca9685"
	"github.com/MarkSaravi/drone-go/utils"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
)

func main() {
	drivers.InitHost()
	flightControlConfigs := config.ReadFlightControlConfig()
	b, _ := i2creg.Open(flightControlConfigs.Configs.ESC.I2CDev)
	i2cConn := &i2c.Dev{Addr: pca9685.PCA9685Address, Bus: b}
	powerbreaker := utils.NewPowerBreaker()
	pca9685.Calibrate(i2cConn, powerbreaker)
}

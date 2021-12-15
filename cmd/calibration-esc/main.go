package main

import (
	"github.com/marksaravi/drone-go/config"
	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/pca9685"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
)

func main() {
	hardware.InitHost()
	configs := config.ReadConfigs().FlightControl.ESC
	b, _ := i2creg.Open(configs.I2CDev)
	i2cConn := &i2c.Dev{Addr: pca9685.PCA9685Address, Bus: b}
	powerbreaker := hardware.NewPowerBreaker()
	pca9685.Calibrate(i2cConn, powerbreaker)
}

package main

import (
	"github.com/MarkSaravi/drone-go/drivers"
	"github.com/MarkSaravi/drone-go/drivers/pca9685"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
)

func main() {
	drivers.InitHost()
	b, _ := i2creg.Open("/dev/i2c-1")
	i2cConn := &i2c.Dev{Addr: pca9685.PCA9685Address, Bus: b}
	powerbreaker := drivers.NewGPIOOutput("GPIO17")
	pca9685.Calibrate(i2cConn, powerbreaker)
}

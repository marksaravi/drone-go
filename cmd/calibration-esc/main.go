package main

import (
	"github.com/marksaravi/drone-go/devices"
	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/pca9685"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
)

func main() {
	hardware.HostInitialize()
	powerBreakerGPIO := hardware.NewGPIOOutput("GPIO17")
	powerBreaker := devices.NewPowerBreaker(powerBreakerGPIO)
	b, _ := i2creg.Open("/dev/i2c-1")
	i2cConn := &i2c.Dev{Addr: pca9685.PCA9685Address, Bus: b}
	pca9685.Calibrate(i2cConn, powerBreaker)
}

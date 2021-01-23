package main

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/devices/pca9685"
	"github.com/MarkSaravi/drone-go/drivers/i2c"
	"github.com/MarkSaravi/drone-go/modules/esc"
	"github.com/MarkSaravi/drone-go/drivers/gpio"
	"github.com/MarkSaravi/drone-go/modules/powerbreaker"
)

func main() {
	i2cConnection, err := i2c.Open("/dev/i2c-1")
	if err != nil {
		fmt.Println(err)
		return
	}
	pca9685, err := pca9685.NewPCA9685Driver(pca9685.PCA9685Address, i2cConnection)
	escs := esc.NewESC(pca9685)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	err = gpio.Open()
	defer gpio.Close()
	pin, err := gpio.NewPin(gpio.GPIO17)
	if err != nil {
		fmt.Println(err)
		return
	}
	breaker := powerbreaker.NewPowerBreaker(pin)
	defer breaker.SetLow()
	defer breaker.SetAsInput()

	escs.Start(float32(esc.Frequency))
	defer escs.Close()
	fmt.Println("setting max pulse width: ", esc.MaxPW)
	fmt.Println("turn on ESCs")
	escs.SetPulseWidthAll(esc.MaxPW)
	time.Sleep(1 * time.Second)
	breaker.SetHigh()
	time.Sleep(12 * time.Second)
	fmt.Println("setting min pulse width: ", esc.MinPW)
	escs.SetPulseWidthAll(esc.MinPW)
	time.Sleep(12 * time.Second)
	fmt.Println("turn off ESCs")
	breaker.SetLow()
	time.Sleep(1 * time.Second)
}

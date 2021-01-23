package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/devices/pca9685"
	"github.com/MarkSaravi/drone-go/drivers/i2c"
	"github.com/MarkSaravi/drone-go/modules/esc"
	"github.com/MarkSaravi/drone-go/drivers/gpio"
	"github.com/MarkSaravi/drone-go/modules/powerbreaker"

)

func main() {
	maxPower := int(25)
	channel := flag.Int("ch", 0, "ESC channel")
	power := flag.Int("power", 7, "Power")
	flag.Parse()
	if *power > maxPower {
		*power = maxPower
	}
	fmt.Println(*channel, esc.Frequency, *power)
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
	defer escs.Close()
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

	escs.Start(esc.Frequency)
	fmt.Println("Starting ", *channel, " at esc.Frequency ", esc.Frequency, " with pulse width ", esc.MinPW)
	escs.SetPulseWidthAll(esc.MinPW)
	time.Sleep(1 * time.Second)
	breaker.SetHigh()
	time.Sleep(5 * time.Second)
	pw := float32(esc.MaxPW - esc.MinPW) * float32(*power) / 100 + float32(esc.MinPW)
	fmt.Println("Starting ", *channel, " at esc.Frequency ", esc.Frequency, " with pulse width ", pw)
	escs.SetPulseWidthAll(pw)
	time.Sleep(5 * time.Second)
	fmt.Println("Starting ", *channel, " at esc.Frequency ", esc.Frequency, " with pulse width ", esc.MinPW)
	escs.SetPulseWidthAll(0.001)
	breaker.SetLow()
	fmt.Println("Stopping...")
}

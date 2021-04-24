package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/devices/pca9685"
	"github.com/MarkSaravi/drone-go/connectors/i2c"
	"github.com/MarkSaravi/drone-go/modules/esc"
	"github.com/MarkSaravi/drone-go/connectors/gpio"
	"github.com/MarkSaravi/drone-go/modules/powerbreaker"

)

func main() {
	channel := flag.Int("ch", 0, "ESC channel")
	flag.Parse()

	maxPower := 25

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
	defer escs.SetPulseWidthAll(0)
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
	fmt.Println("channel: ", *channel, ", PW:  ", esc.MinPW)
	escs.SetPulseWidth(*channel, esc.MinPW)
	breaker.SetHigh()
	time.Sleep(4 * time.Second)
	power := 1
	inc := 1
	for power != 0 {
		pw := float32(esc.MaxPW - esc.MinPW) * float32(power) / 100 + float32(esc.MinPW)
		fmt.Println("channel: ", *channel, ", PW:  ", pw)
		escs.SetPulseWidth(*channel, pw)
		time.Sleep(250 * time.Millisecond)
		power += inc
		if power == maxPower {
			inc = -1
		}
	}
	breaker.SetLow()
	fmt.Println("finished")
}

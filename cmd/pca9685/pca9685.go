package main

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/devices/pca9685"
	"github.com/MarkSaravi/drone-go/drivers/i2c"
)

func main() {
	connection, err := i2c.Open("/dev/i2c-1")
	if err != nil {
		fmt.Println(err)
		return
	}
	pca9685, err := pca9685.NewPCA9685Driver(pca9685.PCA9685Address, connection)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = pca9685.Start(400)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pca9685.Close()
	pca9685.SetPulseWidth(0, 0.002)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Channel 0 PWM is set")
	}
	time.Sleep(1 * time.Second)
}

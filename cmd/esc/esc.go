package main

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/devices/pca9685"
	"github.com/MarkSaravi/drone-go/drivers/i2c"
	"github.com/MarkSaravi/drone-go/modules/esc"
)

func main() {
	i2cConnection, err := i2c.Open("/dev/i2c-1")
	if err != nil {
		fmt.Println(err)
		return
	}
	pca9685, err := pca9685.NewPCA9685Driver(pca9685.PCA9685Address, i2cConnection)
	esc := esc.NewESC(pca9685)
	if err != nil {
		fmt.Println(err)
		return
	}

	esc.Start(350)
	defer esc.Close()
	esc.SetPulseWidth(3, 0.002)
	time.Sleep(5 * time.Second)
}

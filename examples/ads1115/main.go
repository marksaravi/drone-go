package main

import (
	"fmt"

	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/ads1115"
)

func main() {
	hardware.HostInitialize()
	b, _ := i2creg.Open("")
	defer b.Close()
	i2cdev := &i2c.Dev{Addr: 0x48, Bus: b}

	atod := ads1115.NewADS1115(i2cdev);

	for channel:=0; channel<4; channel++ {
		voltage := atod.Read(0)
		fmt.Println(channel, voltage)
	}
	
}
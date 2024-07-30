package main

import (
	"fmt"
	"time"

	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/ads1115"
)

func main() {
	hardware.HostInitialize()
	b, e := i2creg.Open("")
	if e!=nil {
		fmt.Printf(e.Error())
		return
	}
	defer b.Close()
	i2cdev := &i2c.Dev{Addr: 0x48, Bus: b}

	atod := ads1115.NewADS1115(i2cdev);

	for channel:=byte(0); channel<4; channel++ {
		// atod.WriteConfigs(0)
		value:=atod.ReadADC_SingleEnded(channel)
		fmt.Println(value)
		time.Sleep(time.Second/5)
	}
	
}
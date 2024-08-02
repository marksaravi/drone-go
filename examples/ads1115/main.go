package main

import (
	"context"
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

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		fmt.Scanln()
		cancel()
	}()
	atod := ads1115.NewADS1115(i2cdev);

	channel := byte(0)
	running := true
	for running {
		select {
		case <-ctx.Done():
			running = false
		default:
			value:=atod.ReadADC_SingleEnded(channel)
			fmt.Println(value)
		}
		time.Sleep(time.Second/5)
	}
	
}
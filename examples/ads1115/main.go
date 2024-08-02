package main

import (
	"context"
	"fmt"
	"math"
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

	channel := byte(2)
	running := true
	prevValue := float64(0)
	sum:=float64(0)
	average:=float64(0)
	counter:=int(0)
	const SPS = 40
	for running {
		select {
		case <-ctx.Done():
			running = false
		default:
			value:=float64(atod.ReadADC_SingleEnded(channel))
			sum += value
			counter++
			average=sum/float64(counter)

			if math.Abs(prevValue-value)>average/100 {
				fmt.Println()
				fmt.Println(prevValue, value)
			}
			if counter % SPS == 0 {
				fmt.Print(".")
			}
			if counter % SPS*80 == 0 {
				fmt.Print(".")
			}			
			prevValue = value
		}
		time.Sleep(time.Second/SPS)
	}
	
}
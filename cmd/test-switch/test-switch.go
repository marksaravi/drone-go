package main

import (
	"fmt"
	"log"
	"time"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
)

func main() {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}
	var pin gpio.PinIn = gpioreg.ByName("GPIO0")
	pin.In(gpio.PullUp, gpio.NoEdge)
	if pin == nil {
		log.Fatal("Failed to find ")
	}
	var level gpio.Level = gpio.High
	var counter int = 0
	var ptime time.Time = time.Now()
	for {
		nl := pin.Read()
		if nl != level {
			level = nl
			if level == gpio.Low {
				t := time.Now()
				dur := time.Since(ptime)
				ptime = t
				counter++
				if dur > time.Millisecond*100 {
					dur = 0
				}
				fmt.Println(level, counter, dur)
			}
		}
		time.Sleep(time.Millisecond * 10)
	}
}

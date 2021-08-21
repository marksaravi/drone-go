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
	var pin gpio.PinIn = gpioreg.ByName("GPIO26")
	pin.In(gpio.PullDown, gpio.NoEdge)
	if pin == nil {
		log.Fatal("Failed to find ")
	}
	var level gpio.Level = gpio.Low
	var counter int = 0
	var ptime time.Time = time.Now()
	for {
		nl := pin.Read()
		if nl != level {
			level = nl
			if level == gpio.High {
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

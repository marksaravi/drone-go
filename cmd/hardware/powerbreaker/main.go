package main

import (
	"fmt"
	"log"
	"time"

	"github.com/MarkSaravi/drone-go/modules/powerbreaker"
	"periph.io/x/periph/host"
)

func main() {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Started")
	breaker := powerbreaker.NewPowerBreaker("GPIO17")
	breaker.MotorsOn()
	time.Sleep(5 * time.Second)
	breaker.MotorsOff()
}

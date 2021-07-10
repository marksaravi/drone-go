package main

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/modules/powerbreaker"
)

func main() {
	fmt.Println("Started")
	breaker := powerbreaker.NewPowerBreaker()
	breaker.MotorsOn()
	time.Sleep(5 * time.Second)
	breaker.MotorsOff()
}

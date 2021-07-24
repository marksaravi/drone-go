package main

import (
	"fmt"

	"github.com/MarkSaravi/drone-go/hardware/nrf204"
)

func main() {
	fmt.Println("Start")
	receiver := nrf204.CreateNRF204()
	receiver.Init()
	receiver.OpenReadingPipe()
	receiver.SetPALevel()
	receiver.StartListening()
	receiver.IsAvailable()
	data := receiver.Read()
	fmt.Println(data)
	fmt.Println("End")
}

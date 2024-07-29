package main

import (
	"fmt"
	"github.com/marksaravi/drone-go/hardware/ads1115"
)

func main() {
	atod := ads1115.NewADS1115();

	for channel:=0; channel<4; channel++ {
		voltage := atod.Read(0)
		fmt.Println(channel, voltage)
	}
	
}
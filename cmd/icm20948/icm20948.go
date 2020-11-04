package main

import (
	"fmt"

	"github.com/MarkSaravi/drone-go/modules/mpu"
)

func main() {
	mpu, err := mpu.NewMPU("/dev/spidev0.0")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer mpu.Close()
	if whoami, err := mpu.WhoAmI(); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("Who am I: 0x%X\n", whoami)
	}
}

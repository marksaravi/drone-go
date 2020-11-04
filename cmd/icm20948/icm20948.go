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
	if err = mpu.SetFullScaleRange(0b00000000); err != nil {
		fmt.Println(err.Error())
	}
	if fsr, err := mpu.GetFullScaleRange(); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("Gyro full scale: 0x%X\n", fsr)
	}
}

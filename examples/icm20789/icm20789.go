package main

import (
	"context"
	"fmt"
	"log"

	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/icm20789"
	"github.com/marksaravi/drone-go/timeinterval"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	hardware.HostInitialize()
	icm20789Configs := icm20789.ReadConfigs("./configs/hardware.json")
	ctx, cancel := context.WithCancel(context.Background())

	go func(cancel context.CancelFunc) {
		fmt.Scanln()
		cancel()
	}(cancel)

	mems := icm20789.NewICM20789(icm20789Configs)
	whoAmI, err := mems.WhoAmI()
	if err == nil {
		fmt.Printf("WHO AM I: %x\n", whoAmI)
	}

	readInterval := timeinterval.WithDataPerSecond(1000)

	running := true
	for running {
		select {
		case <-ctx.Done():
			running = false
		default:
			if readInterval.IsTime() {
				data, _ := mems.Read()
				acc := data.Accelerometer
				gyro := data.Gyroscope

				log.Printf("Accelerometer  X: %6.2f, Y: %6.2f, Gyro  X: %6.2f, Y: %6.2f", acc.X, acc.Y, gyro.DX, gyro.DY)
			}
		}

	}
}

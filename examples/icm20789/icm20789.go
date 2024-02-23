package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/icm20789"
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

	lastRead := time.Now()
	var maxX float64 = math.SmallestNonzeroFloat64
	maxY := maxX
	var minX float64 = math.MaxFloat64
	minY := minX
	running := true
	counter := 0
	for running {
		select {
		case <-ctx.Done():
			running = false
		default:
			if time.Since(lastRead) >= time.Second {
				lastRead = time.Now()
				data, _ := mems.Read()
				acc := data.Accelerometer
				// gyro:=data.Gyroscope

				counter++
				if counter > 1 {
					if acc.X > maxX {
						maxX = acc.X
					}
					if acc.X < minX {
						minX = acc.X
					}
					if acc.Y > maxY {
						maxY = acc.Y
					}
					if acc.Y < minY {
						minY = acc.Y
					}
				}
				log.Printf("Accelerometer  X: %6.2f, Y: %6.2f", acc.X, acc.Y)
				// if time.Since(lastPrint) >= time.Second/2 {
				// 	log.Printf("Accelerometer  X: %6.2f, Y: %6.2f", acc.X, acc.Y)
				// 	lastPrint=time.Now()
				// }
			}
		}

	}
	fmt.Printf("Acc dX:  %0.2f, minX:  %0.2f, maxX:  %0.2f\n\n", maxX-minX, minX, maxX)
	fmt.Printf("Acc dY:  %0.2f, minY:  %0.2f, maxY:  %0.2f\n\n", maxY-minY, minY, maxY)
}

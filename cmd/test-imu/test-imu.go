package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/marksaravi/drone-go/devices/imu"
	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/mems/icm20789"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	hardware.HostInitialize()
	ctx, cancel := context.WithCancel(context.Background())

	go func(cancel context.CancelFunc) {
		fmt.Scanln()
		cancel()
	}(cancel)

	mems := icm20789.NewICM20789()

	imudev := imu.NewIMU(mems)

	lastRead := time.Now()
	ticker := time.NewTicker(time.Second / 2)
	var rotations imu.Rotations
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			log.Printf("%6.2f, %6.2f, %6.2f\n", rotations.Roll, rotations.Pitch, rotations.Yaw)
		default:
			if time.Since(lastRead) >= time.Millisecond {
				lastRead = time.Now()
				rotations, _ = imudev.Read()
			}
		}

	}
}

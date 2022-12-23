package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/icm20789"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	hardware.HostInitialize()
	ctx, cancel := context.WithCancel(context.Background())

	go func(cancel context.CancelFunc) {
		fmt.Scanln()
		cancel()
	}(cancel)

	imu := icm20789.NewICM20789()
	imu.Setup()

	lastRead := time.Now()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if time.Since(lastRead) > time.Second/2 {
				lastRead = time.Now()
				data, err := imu.ReadIMUData()
				if err == nil {
					log.Println(data)
				} else {
					log.Println("IVALID DATA: ", err)
				}
			}
		}

	}
}

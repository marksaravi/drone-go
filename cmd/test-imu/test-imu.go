package main

import (
	"context"
	"fmt"
	"time"

	"github.com/marksaravi/drone-go/hardware/icm20789"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go func(cancel context.CancelFunc) {
		fmt.Scanln()
		cancel()
	}(cancel)

	imu := icm20789.NewIcm20987(0, 0)
	imu.Initialize()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		data, err := imu.ReadAccelerometer()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(data)
		time.Sleep(time.Second)
	}
}

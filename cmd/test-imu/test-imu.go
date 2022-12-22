package main

import (
	"fmt"

	"github.com/marksaravi/drone-go/hardware/icm20789"
)

func main() {
	// ctx, cancel := context.WithCancel(context.Background())

	// go func(cancel context.CancelFunc) {
	// 	fmt.Scanln()
	// 	cancel()
	// }(cancel)

	imu := icm20789.NewIcm20987(0, 0)
	whoami, powmgm1, readok, readerr := imu.SPIReadTest()
	if readok {
		fmt.Printf("WHO AM I: 0x%x, POWER MANAGEMENT 1: 0x%x\n", whoami, powmgm1)
	} else {
		fmt.Printf("spi error: %v\n", readerr)
	}
	writeok, writeerr := imu.SPIWriteTest()
	if writeok {
		fmt.Println("Write ok")
	} else {
		fmt.Println(writeerr)
	}

	// imu.Initialize("4g")
	// lastReadTime := time.Now()
	// for {
	// 	select {
	// 	case <-ctx.Done():
	// 		return
	// 	default:
	// 		if time.Since(lastReadTime) >= time.Second {
	// 			data, err := imu.ReadAccelerometer()
	// 			if err != nil {
	// 				fmt.Println(err)
	// 			}
	// 			fmt.Println(data)
	// 			lastReadTime = time.Now()
	// 		}
	// 	}
	// }
}

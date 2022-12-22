package main

import (
	"log"
	"time"

	"github.com/marksaravi/drone-go/hardware/icm20789"
)

func main() {
	// ctx, cancel := context.WithCancel(context.Background())

	// go func(cancel context.CancelFunc) {
	// 	fmt.Scanln()
	// 	cancel()
	// }(cancel)

	imu := icm20789.NewIcm20987(0, 0)
	whoami, _ := imu.WhoAmI()
	log.Println("WHOAMI: ", whoami)
	imu.Initialize("4g")
	time.Sleep(100 * time.Millisecond)
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

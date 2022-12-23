package main

import (
	"log"

	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/icm20789"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	hardware.HostInitialize()
	// ctx, cancel := context.WithCancel(context.Background())

	// go func(cancel context.CancelFunc) {
	// 	fmt.Scanln()
	// 	cancel()
	// }(cancel)

	imu := icm20789.NewICM20789()
	imu.Setup()
}

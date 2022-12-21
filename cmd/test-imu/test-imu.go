package main

import (
	"fmt"
	"time"

	"github.com/marksaravi/drone-go/hardware/icm20789"
)

func main() {
	imu := icm20789.NewIcm20987(0, 0)
	for {
		data, err := imu.WhoAmI()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(data)
		time.Sleep(time.Second)
	}
}

package main

import (
	"fmt"

	"github.com/MarkSaravi/drone-go/devices/icm20948"
)

func main() {
	driver, _ := icm20948.NewRaspberryPiDriver(0, 0)
	fmt.Println(driver)
}

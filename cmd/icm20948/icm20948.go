package main

import (
	"fmt"

	"github.com/MarkSaravi/drone-go/devices/icm20948"
)

func main() {
	// b, err := sysfs.NewSPI(0, 0)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer b.Close()

	// c, err := b.Connect(7*physic.MegaHertz, spi.Mode3, 8)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// r := make([]byte, 2)
	// if err := c.Tx([]byte{0b10000000, 0x0}, r); err != nil {
	// 	fmt.Println(err.Error())
	// }
	// fmt.Printf("%X\n", r)
	driver, err := icm20948.NewRaspberryPiICM20948Driver("/dev/spidev0.0")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if data, err := driver.Read(icm20948.WHO_AM_I); err == nil {
		fmt.Printf("0x%X\n", data)
	} else {
		fmt.Println(err.Error())
	}
	defer driver.Close()
}

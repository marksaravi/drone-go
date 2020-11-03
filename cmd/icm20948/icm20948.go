package main

import (
	"fmt"
	"log"

	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/host/sysfs"
)

func main() {
	// icm20948, err := icm20948.NewRaspberryPiDriver(0, 0)
	// defer icm20948.Close()
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// fmt.Println(icm20948)
	// tx := make([]byte, 2)
	// tx[0] = 0
	// tx[1] = 0
	// rx := make([]byte, 2)
	// err = icm20948.Connection.Tx([]byte{0x1, 0x0}, rx)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println(rx)
	// }
	b, err := sysfs.NewSPI(0, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer b.Close()

	c, err := b.Connect(7*physic.MegaHertz, spi.Mode3, 8)

	if err != nil {
		log.Fatal(err)
	}

	r := make([]byte, 2)
	if err := c.Tx([]byte{0b10000000, 0x0}, r); err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("%X\n", r)
}

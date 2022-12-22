package main

import (
	"fmt"
	"log"

	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/host/v3/sysfs"
)

type imuIcm20789 struct {
	spiConn spi.Conn
}

func (imu *imuIcm20789) readRegister(address byte, size int) ([]byte, error) {
	w := make([]byte, size+1)
	r := make([]byte, size+1)
	w[0] = address | byte(0x80)

	err := imu.spiConn.Tx(w, r)
	return r[1:], err
}

func (imu *imuIcm20789) readByteFromRegister(address byte) (byte, error) {
	res, err := imu.readRegister(address, 1)
	return res[0], err
}

func (imu *imuIcm20789) writeRegister(address byte, data ...byte) error {
	w := make([]byte, 1, len(data)+1)
	r := make([]byte, cap(w))
	w[0] = address
	w = append(w, data...)
	fmt.Println("len: ", len(w), len(r))
	err := imu.spiConn.Tx(w, r)
	return err
}

func main() {
	initialize()
	// Make sure periph is initialized.
	// TODO: Use host.Init(). It is not used in this example to prevent circular
	// go package import.
	// if _, err := driverreg.Init(); err != nil {
	// 	log.Fatal(err)
	// }
	fmt.Println("initializing SPI")

	// Use spireg SPI port registry to find the first available SPI bus.
	// p, err := spireg.Open("")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	p, err := sysfs.NewSPI(0, 0)

	if err != nil {
		log.Fatal(err)
	}

	// Convert the spi.Port into a spi.Conn so it can be used for communication.
	c, err := p.Connect(physic.MegaHertz, spi.Mode3, 8)

	if err != nil {
		log.Fatal(err)
	}

	// Write 0x10 to the device, and read a byte right after.
	write := []byte{0x10, 0x00}
	read := make([]byte, len(write))
	if err := c.Tx(write, read); err != nil {
		log.Fatal(err)
	}
	// Use read.
	fmt.Printf("READ VALUE%v\n", read[1:])
}

func NewSPIConnection(busNumber int, chipSelect int) spi.Conn {
	p, err := sysfs.NewSPI(0, 0)

	if err != nil {
		log.Fatal(err)
	}

	// Convert the spi.Port into a spi.Conn so it can be used for communication.
	c, err := p.Connect(physic.MegaHertz, spi.Mode3, 8)

	if err != nil {
		log.Fatal(err)
	}
	return c
}

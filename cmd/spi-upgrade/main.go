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
	w[0] = address
	w = append(w, data...)
	err := imu.spiConn.Tx(w, nil)
	return err
}

func main() {
	initialize()
	fmt.Println("initializing SPI")
	c := NewSPIConnection(0, 0)
	imu := imuIcm20789{
		spiConn: c,
	}

	whoami, _ := imu.readByteFromRegister(0x75)
	fmt.Printf("WHO AM I: 0x%x\n", whoami)

	power, _ := imu.readByteFromRegister(107)
	fmt.Printf("POWER: 0x%x\n", power)

	const WRITE_ADDRESS = 19
	const NEW_VALUE byte = 13
	originalValue, _ := imu.readByteFromRegister(WRITE_ADDRESS)
	fmt.Printf("ORIGINAL: 0x%x\n", originalValue)
	imu.writeRegister(WRITE_ADDRESS, NEW_VALUE)
	newValue, err := imu.readByteFromRegister(WRITE_ADDRESS)
	fmt.Printf("NEW: 0x%x, %v\n", newValue, err)
}

func NewSPIConnection(busNumber int, chipSelect int) spi.Conn {
	p, err := sysfs.NewSPI(0, 0)

	if err != nil {
		log.Fatal(err)
	}

	// Convert the spi.Port into a spi.Conn so it can be used for communication.
	c, err := p.Connect(physic.MegaHertz, spi.Mode0, 8)

	if err != nil {
		log.Fatal(err)
	}
	return c
}

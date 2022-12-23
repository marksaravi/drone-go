package main

import (
	"fmt"

	"github.com/marksaravi/drone-go/hardware2"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
)

type imuIcm20789 struct {
	i2cConn *i2c.Dev
}

func (imu *imuIcm20789) readRegister(address byte, size int) ([]byte, error) {
	return []byte{}, nil
}

func (imu *imuIcm20789) writeRegister(address byte, data ...byte) error {
	return nil
}

func (imu *imuIcm20789) readByteFromRegister(address byte) (byte, error) {
	res, err := imu.readRegister(address, 1)
	return res[0], err
}

func main() {
	hardware2.InitializeHost()
	b, err := i2creg.Open("/dev/i2c-1")
	if err != nil {
		fmt.Println("I2CREG ERROR: ", err)
		return
	}
	defer b.Close()
	const ADDRESS uint16 = 0b1101000

	// Dev is a valid conn.Conn.
	d := &i2c.Dev{Addr: ADDRESS, Bus: b}

	// Send a command 0x10 and expect a 5 bytes reply.
	const command byte = 23
	write := []byte{command}
	read := make([]byte, 2)
	if err := d.Tx(write, read); err != nil {
		fmt.Println("READ ERROR: ", err)
		return
	}
	fmt.Printf("%v\n", read)
}

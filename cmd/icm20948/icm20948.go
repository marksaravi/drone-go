package main

import (
	"fmt"
	"os"

	"github.com/MarkSaravi/drone-go/devices/icm20948"

	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/host/sysfs"
)

type IMUDevice struct {
	*sysfs.SPI
	spi.Conn
	regbank byte
}

func (dev *IMUDevice) SelRegisterBank(regbank byte) error {
	if regbank == dev.regbank {
		return nil
	}
	dev.regbank = regbank

	fmt.Printf("SelRegisterBank to %d\n", dev.regbank)
	return dev.Conn.Tx([]byte{icm20948.REG_BANK_SEL, (regbank << 4) & 0x30}, nil)
}

func (dev *IMUDevice) ReadRegister(address byte, len int) ([]byte, error) {
	w := make([]byte, len+1)
	r := make([]byte, len+1)
	w[0] = (address & 0x7F) | 0x80
	// defer prn(fmt.Sprintf("ReadRegister (0x%X)", address), r)
	err := dev.Conn.Tx(w, r)
	return r[1:], err
}

func (dev *IMUDevice) WriteRegister(address byte, data ...byte) error {
	// defer prn(fmt.Sprintf("ReadRegister (0x%X)", address), r)
	if len(data) == 0 {
		return nil
	}
	w := append([]byte{address & 0x7F}, data...)
	err := dev.Conn.Tx(w, nil)
	return err
}

func errCheck(step string, err error) {
	if err != nil {
		fmt.Printf("Error at %s: %s\n", step, err.Error())
		os.Exit(0)
	}
}

func prn(msg string, bytes []byte) {
	fmt.Printf("%s: ", msg)
	for _, b := range bytes {
		fmt.Printf("0x%X, ", b)
	}
	fmt.Printf("\n")
}

func main() {
	r := make([]byte, 2)

	d, err := sysfs.NewSPI(0, 0)
	errCheck("sysfs.NewSPI", err)
	conn, err := d.Connect(7*physic.MegaHertz, spi.Mode3, 8)
	errCheck("spi.Connect", err)
	dev := &IMUDevice{
		SPI:     d,
		Conn:    conn,
		regbank: 0xFF,
	}

	errCheck("SelRegisterBank", dev.SelRegisterBank(0))

	r, err = dev.ReadRegister(icm20948.WHO_AM_I, 1)
	prn("Who am I", r)

	// set bank 2
	// dev.Conn.Tx([]byte{icm20948.REG_BANK_SEL, icm20948.BANK2}, nil)
	errCheck("SelRegisterBank", dev.SelRegisterBank(2))

	// read MOD_CTRL_USR
	r, err = dev.ReadRegister(icm20948.MOD_CTRL_USR, 1)
	prn("MOD_CTRL_USR bank2", r)

	r, err = dev.ReadRegister(icm20948.WHO_AM_I, 1)
	prn("Who am I", r)

	// read PWR_MGMT_1
	errCheck("SelRegisterBank", dev.SelRegisterBank(0))
	r, err = dev.ReadRegister(icm20948.PWR_MGMT_1, 1)
	prn("PWR_MGMT_1 bank0", r)

	r, err = dev.ReadRegister(icm20948.WHO_AM_I, 1)
	prn("Who am I", r)

	// read PWR_MGMT_1
	errCheck("SelRegisterBank", dev.SelRegisterBank(0))
	err = dev.WriteRegister(icm20948.PWR_MGMT_2, 0b00000111)
	r, err = dev.ReadRegister(icm20948.PWR_MGMT_2, 1)
	prn("PWR_MGMT_2 bank0", r)
	err = dev.WriteRegister(icm20948.PWR_MGMT_2, 0b00111000)
	r, err = dev.ReadRegister(icm20948.PWR_MGMT_2, 1)
	prn("PWR_MGMT_2 bank0", r)

	errCheck("SelRegisterBank", dev.SelRegisterBank(2))
	err = dev.WriteRegister(icm20948.GYRO_SMPLRT_DIV, 2)
	r, err = dev.ReadRegister(icm20948.GYRO_SMPLRT_DIV, 1)
	prn("GYRO_SMPLRT_DIV", r)
}

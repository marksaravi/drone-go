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

func (dev *IMUDevice) SelRegisterBank(bank byte) error {
	var regbank byte = (bank << 4) & 0x30
	if regbank == dev.regbank {
		return nil
	}
	dev.regbank = regbank

	fmt.Printf("SelRegisterBank to %d\n", bank)
	return dev.Conn.Tx([]byte{icm20948.REG_BANK_SEL, regbank}, nil)
}

func (dev *IMUDevice) ReadRegister(address byte) ([]byte, error) {
	r := make([]byte, 2)
	// defer prn(fmt.Sprintf("ReadRegister (0x%X)", address), r)
	err := dev.Conn.Tx([]byte{(address & 0x7F) | 0x80, 0}, r)
	return r, err
}

func errCheck(step string, err error) {
	if err != nil {
		fmt.Printf("Error at %s: %s\n", step, err.Error())
		os.Exit(0)
	}
}

func prn(msg string, b []byte) {
	fmt.Printf("%s: 0x%X, 0x%X\n---------------------\n", msg, b[0], b[1])
}

func main() {
	r := make([]byte, 2)

	d, err := sysfs.NewSPI(0, 0)
	errCheck("sysfs.NewSPI", err)
	conn, err := d.Connect(7*physic.MegaHertz, spi.Mode3, 8)
	errCheck("spi.Connect", err)
	dev := &IMUDevice{
		SPI:  d,
		Conn: conn,
	}

	errCheck("SelRegisterBank", dev.SelRegisterBank(0))

	r, err = dev.ReadRegister(icm20948.WHO_AM_I)
	prn("Who am I", r)

	// set bank 2
	// dev.Conn.Tx([]byte{icm20948.REG_BANK_SEL, icm20948.BANK2}, nil)
	errCheck("SelRegisterBank", dev.SelRegisterBank(2))

	// read MOD_CTRL_USR
	r, err = dev.ReadRegister(icm20948.MOD_CTRL_USR)
	prn("MOD_CTRL_USR bank2", r)

	r, err = dev.ReadRegister(icm20948.WHO_AM_I)
	prn("Who am I", r)

	// read PWR_MGMT_1
	errCheck("SelRegisterBank", dev.SelRegisterBank(0))
	r, err = dev.ReadRegister(icm20948.PWR_MGMT_1)
	prn("PWR_MGMT_1 bank0", r)

	r, err = dev.ReadRegister(icm20948.WHO_AM_I)
	prn("Who am I", r)

	// read PWR_MGMT_1
	errCheck("SelRegisterBank", dev.SelRegisterBank(0))
	r, err = dev.ReadRegister(icm20948.PWR_MGMT_2)
	prn("PWR_MGMT_2 bank0", r)
}

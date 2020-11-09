package icm20948

import (
	"fmt"
	"os"

	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/host"
	"periph.io/x/periph/host/sysfs"
)

const (
	REG_BANK_SEL byte = 0x7F

	// BANK0
	WHO_AM_I     byte = 0x0
	LP_CONFIG    byte = 0x5
	PWR_MGMT_1   byte = 0x6
	PWR_MGMT_2   byte = 0x7
	INT_ENABLE_3 byte = 0x13
	ACCEL_ZOUT_H byte = 0x31
	ACCEL_ZOUT_L byte = 0x32
	GYRO_ZOUT_L  byte = 0x38

	// BANK1
	XA_OFFS_H byte = 0x14

	// BANK2
	GYRO_SMPLRT_DIV byte = 0x0
	GYRO_CONFIG_1   byte = 0x1
	GYRO_CONFIG_2   byte = 0x2
	ZG_OFFS_USRL    byte = 0x8
	ACCEL_CONFIG_2  byte = 0x15
	MOD_CTRL_USR    byte = 0x54
)

type IMUDevice struct {
	*sysfs.SPI
	spi.Conn
	regbank byte
}

func init() {
	host.Init()
}

// NewRaspberryPiICM20948Driver creates ICM20948 driver for raspberry pi
func NewRaspberryPiICM20948Driver(busNumber int, chipSelect int) (*IMUDevice, error) {
	d, err := sysfs.NewSPI(busNumber, chipSelect)
	if err != nil {
		return nil, err
	}
	conn, err := d.Connect(7*physic.MegaHertz, spi.Mode3, 8)
	if err != nil {
		return nil, err
	}
	dev := IMUDevice{
		SPI:     d,
		Conn:    conn,
		regbank: 0xFF,
	}
	return &dev, nil
}

func (dev *IMUDevice) SelRegisterBank(regbank byte) error {
	if regbank == dev.regbank {
		return nil
	}
	dev.regbank = regbank

	fmt.Printf("SelRegisterBank to %d\n", dev.regbank)
	return dev.WriteRegister(REG_BANK_SEL, (regbank<<4)&0x30)
}

func (dev *IMUDevice) ReadRegister(address byte, len int) ([]byte, error) {
	w := make([]byte, len+1)
	r := make([]byte, len+1)
	w[0] = (address & 0x7F) | 0x80
	// defer Prn(fmt.Sprintf("ReadRegister (0x%X)", address), r)
	fmt.Println(r)
	err := dev.Conn.Tx(w, r)
	return r[1:], err
}

func (dev *IMUDevice) WriteRegister(address byte, data ...byte) error {
	// defer Prn(fmt.Sprintf("ReadRegister (0x%X)", address), r)
	if len(data) == 0 {
		return nil
	}
	w := append([]byte{address & 0x7F}, data...)
	err := dev.Conn.Tx(w, nil)
	return err
}

func ErrCheck(step string, err error) {
	if err != nil {
		fmt.Printf("Error at %s: %s\n", step, err.Error())
		os.Exit(0)
	}
}

func Prn(msg string, bytes []byte) {
	fmt.Printf("%s: ", msg)
	for _, b := range bytes {
		fmt.Printf("0x%X, ", b)
	}
	fmt.Printf("\n")
}

package icm20948

import (
	"fmt"

	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/host"
	"periph.io/x/periph/host/sysfs"
)

const (
	BANK0 uint16 = 0 << 8
	BANK1 uint16 = 1 << 8
	BANK2 uint16 = 2 << 8
	BANK3 uint16 = 3 << 8
)

const (
	REG_BANK_SEL byte = 0x7F

	// BANK0
	WHO_AM_I     uint16 = BANK0 | 0x0
	LP_CONFIG    uint16 = BANK0 | 0x5
	PWR_MGMT_1   uint16 = BANK0 | 0x6
	PWR_MGMT_2   uint16 = BANK0 | 0x7
	INT_ENABLE_3 uint16 = BANK0 | 0x13
	ACCEL_ZOUT_H uint16 = BANK0 | 0x31
	ACCEL_ZOUT_L uint16 = BANK0 | 0x32
	GYRO_ZOUT_L  uint16 = BANK0 | 0x38

	// BANK1
	XA_OFFS_H uint16 = BANK1 | 0x14

	// BANK2
	GYRO_SMPLRT_DIV uint16 = BANK2 | 0x0
	GYRO_CONFIG_1   uint16 = BANK2 | 0x1
	GYRO_CONFIG_2   uint16 = BANK2 | 0x2
	ZG_OFFS_USRL    uint16 = BANK2 | 0x8
	ACCEL_CONFIG_2  uint16 = BANK2 | 0x15
	MOD_CTRL_USR    uint16 = BANK2 | 0x54
)

func reg(reg uint16) *Register {
	return &Register{
		address: byte(reg),
		bank:    byte(reg >> 8),
	}
}

type Register struct {
	address byte
	bank    byte
}

type Device struct {
	*sysfs.SPI
	spi.Conn
	regbank byte
}

func init() {
	host.Init()
}

// NewRaspberryPiICM20948Driver creates ICM20948 driver for raspberry pi
func NewRaspberryPiICM20948Driver(busNumber int, chipSelect int) (*Device, error) {
	d, err := sysfs.NewSPI(busNumber, chipSelect)
	if err != nil {
		return nil, err
	}
	conn, err := d.Connect(7*physic.MegaHertz, spi.Mode3, 8)
	if err != nil {
		return nil, err
	}
	dev := Device{
		SPI:     d,
		Conn:    conn,
		regbank: 0xFF,
	}
	return &dev, nil
}

func (dev *Device) readReg(address byte, len int) ([]byte, error) {
	w := make([]byte, len+1)
	r := make([]byte, len+1)
	w[0] = (address & 0x7F) | 0x80
	err := dev.Conn.Tx(w, r)
	return r[1:], err
}

func (dev *Device) writeReg(address byte, data ...byte) error {
	if len(data) == 0 {
		return nil
	}
	w := append([]byte{address & 0x7F}, data...)
	err := dev.Conn.Tx(w, nil)
	return err
}

func (dev *Device) selRegisterBank(regbank byte) error {
	if regbank == dev.regbank {
		return nil
	}
	dev.regbank = regbank
	fmt.Println("Switch to bank", dev.regbank)
	return dev.writeReg(REG_BANK_SEL, (regbank<<4)&0x30)
}

func (dev *Device) ReadRegister(register uint16, len int) ([]byte, error) {
	reg := reg(register)
	dev.selRegisterBank(reg.bank)
	return dev.readReg(reg.address, len)
}

func (dev *Device) WriteRegister(register uint16, data ...byte) error {
	if len(data) == 0 {
		return nil
	}
	reg := reg(register)
	dev.selRegisterBank(reg.bank)
	return dev.writeReg(reg.address, data...)
}

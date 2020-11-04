package icm20948

import (
	"fmt"

	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/conn/spi/spireg"
	"periph.io/x/periph/host"
)

const (
	BANK0 byte = 0b00000000
	BANK1 byte = 0b00010000
	BANK2 byte = 0b00100000
	BANK3 byte = 0b00110000
)

const (
	WHO_AM_I      byte = 0x0
	GYRO_CONFIG_1 byte = 0x1
	REG_BANK_SEL  byte = 0x7F
)

// Driver for ICM20948
type Driver struct {
	name string
	spi.PortCloser
	spi.Conn
}

func init() {
	host.Init()
}

// NewRaspberryPiICM20948Driver creates ICM20948 driver for raspberry pi
func NewRaspberryPiICM20948Driver(devname string) (*Driver, error) {
	dev, err := spireg.Open(devname)
	if err != nil {
		return nil, err
	}
	connection, err := dev.Connect(7*physic.MegaHertz, spi.Mode3, 8)
	if err != nil {
		return nil, err
	}
	return &Driver{
		name:       "ICM-20948",
		PortCloser: dev,
		Conn:       connection,
	}, nil
}

//Close closes the device
func (d *Driver) Close() {
	d.PortCloser.Close()
}

func (d *Driver) read(address byte) (byte, error) {
	r := make([]byte, 2)
	err := d.Conn.Tx([]byte{0b10000000 | address, 0x0}, r)
	fmt.Printf("raw data: %v\n", r)
	return r[1], err
}

func (d *Driver) write(address byte, data byte) error {
	r := make([]byte, 2)
	return d.Conn.Tx([]byte{address, data}, r)
}

func (d *Driver) selRegisterBank(bank byte) error {
	return d.write(REG_BANK_SEL, bank)
}

func (d *Driver) writeRegister(address byte, bank byte, data byte) error {
	if err := d.selRegisterBank(bank); err != nil {
		return err
	}
	return d.write(address, data)
}

func (d *Driver) readRegister(address byte, bank byte) (byte, error) {
	if err := d.selRegisterBank(bank); err != nil {
		return 0, err
	}
	return d.read(address)
}

// WhoAmI is reading the device
func (d *Driver) WhoAmI() (byte, error) {
	return d.readRegister(WHO_AM_I, BANK0)
}

// SetFullScaleRange to setup Gyroscope range
func (d *Driver) SetFullScaleRange(data byte) error {
	return d.writeRegister(GYRO_CONFIG_1, BANK2, data)
}

// GetFullScaleRange to read Gyroscope range
func (d *Driver) GetFullScaleRange() (byte, error) {
	return d.readRegister(GYRO_CONFIG_1, BANK2)
}

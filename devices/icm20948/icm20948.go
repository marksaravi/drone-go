package icm20948

import (
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
	// BANK0
	WHO_AM_I     byte = 0x0
	REG_BANK_SEL byte = 0x7F

	// BANK2
	GYRO_SMPLRT_DIV byte = 0x0
	GYRO_CONFIG_1   byte = 0x1
	MOD_CTRL_USR    byte = 0x54
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
	connection, err := dev.Connect(7*physic.MegaHertz, spi.Mode0, 8)
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

func (d *Driver) write(address, data byte) error {
	r := []byte{0, 0}
	err := d.Conn.Tx([]byte{address, data}, r)
	// fmt.Printf("write address: %d, data: 0x%0x\n", address, data)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	return err
}

func (d *Driver) read(address byte) ([]byte, error) {
	r := []byte{0, 0}
	err := d.Conn.Tx([]byte{0b10000000 | address, 0}, r)
	// fmt.Printf("read  address: %d, b0: 0x%0x, b1: 0x%X\n", address, r[0], r[1])
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	return r, err
}

func (d *Driver) selRegisterBank(bank byte) error {
	return d.write(0x7F, bank)
}

// SetRegister to setup Gyroscope range
func (d *Driver) SetRegister(address, bank, data byte) error {
	err := d.selRegisterBank(bank)
	if err != nil {
		return err
	}
	return d.write(address, data)
}

// GetRegister to read Gyroscope range
func (d *Driver) GetRegister(address, bank byte) (byte, error) {
	r := []byte{0, 0}
	err := d.selRegisterBank(bank)
	if err != nil {
		return 0, err
	}
	r, err = d.read(address)
	return r[1], err
}

package icm20948

import (
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/conn/spi/spireg"
	"periph.io/x/periph/host"
)

const (
	WHO_AM_I = 0x0
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
	return r[1], err
}

// WhoAmI is reading the device
func (d *Driver) WhoAmI() (byte, error) {
	return d.read(WHO_AM_I)
}

// SetFullScaleRange to setup Gyroscope range
func (d *Driver) SetFullScaleRange() {
}

// GetFullScaleRange to read Gyroscope range
func (d *Driver) GetFullScaleRange() {
}

package icm20948

import (
	"github.com/MarkSaravi/drone-go/drivers/spi"
	xspi "periph.io/x/periph/conn/spi"
)

type Driver struct {
	name string
	spi.Connection
}

// NewDriver creates new ICM20948 driver
func NewDriver(busNum int, chipNum int, mode xspi.Mode, maxSpeed int64, bits int) (*Driver, error) {
	connection, err := spi.GetSpiConnection(busNum, chipNum, mode, maxSpeed, bits)
	if err != nil {
		return nil, err
	}
	return &Driver{
		name:       "ICM-20948",
		Connection: connection,
	}, nil
}

// NewRaspberryPiDriver for raspberry pi
func NewRaspberryPiDriver(busNum int, chipNum int) (*Driver, error) {
	return NewDriver(busNum, chipNum, xspi.Mode3, 7000000, 128)
}

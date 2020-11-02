package icm20948

import (
	"github.com/MarkSaravi/drone-go/drivers/spi"
	xspi "periph.io/x/periph/conn/spi"
)

type Driver struct {
	name string
	spi.Connection
}

func NewDriver() (*Driver, error) {
	connection, err := spi.GetSpiConnection(0, 0, xspi.Mode3, 7000000, 128)
	if err != nil {
		return nil, err
	}
	return &Driver{
		name:       "ICM-20948",
		Connection: connection,
	}, nil
}

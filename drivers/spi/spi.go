package spi

import (
	"periph.io/x/periph/conn/physic"
	xspi "periph.io/x/periph/conn/spi"
	xsysfs "periph.io/x/periph/host/sysfs"
)

// Operations are the wrappers around the actual functions used by the SPI device interface
type Operations interface {
	Close() error
	Tx(w, r []byte) error
}

/// Connection is a connection to a SPI device with a specific bus/chip.
// Provided by an Adaptor, usually just by calling the spi package's GetSpiConnection() function.
type Connection Operations

// SpiConnection is the implementation of the SPI interface using the periph.io
// sysfs implementation for Linux.
type SpiConnection struct {
	Operations
	port     xspi.PortCloser
	dev      xspi.Conn
	bus      int
	chip     int
	bits     int
	mode     int
	maxSpeed int64
}

// NewConnection creates and returns a new connection to a specific
// spi device on a bus/chip using the periph.io interface.
func NewConnection(port xspi.PortCloser, conn xspi.Conn) (connection *SpiConnection) {
	return &SpiConnection{port: port, dev: conn}
}

// Close the SPI connection.
func (c *SpiConnection) Close() error {
	return c.port.Close()
}

// Tx uses the SPI device to send/receive data.
func (c *SpiConnection) Tx(w, r []byte) error {
	return c.dev.Tx(w, r)
}

// GetSpiConnection is a helper to return a SPI device.
func GetSpiConnection(busNum int, chipNum int, mode xspi.Mode, maxSpeed int64, bits int) (Connection, error) {
	p, err := xsysfs.NewSPI(busNum, chipNum)
	if err != nil {
		return nil, err
	}
	c, err := p.Connect(physic.Frequency(maxSpeed)*physic.Hertz, xspi.Mode(mode), bits)
	if err != nil {
		return nil, err
	}
	return NewConnection(p, c), nil
}

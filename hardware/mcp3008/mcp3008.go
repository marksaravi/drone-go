package mcp3008

import (
	"periph.io/x/periph/conn/spi"
)

type mcp3008dev struct {
	spiConn      spi.Conn
	channel      int
	valueRange   int
	digitalRange int
	midValue     int
	value        int
}

func NewMCP3008(spiConn spi.Conn, channel int, valueRange int, digitalRange int, midValue int) *mcp3008dev {
	return &mcp3008dev{
		spiConn:      spiConn,
		channel:      channel,
		valueRange:   valueRange,
		digitalRange: digitalRange,
		midValue:     midValue,
	}
}

func (dev *mcp3008dev) Read() int {
	ch := byte(dev.channel)
	if ch > 7 {
		ch = 0
	}
	w := []byte{0b00000001, 0b10000000 | (ch << 4), 0b00000000}
	r := []byte{0, 0, 0}
	err := dev.spiConn.Tx(w, r)
	var digitalValue uint16 = uint16(r[2]) | (uint16(r[1]) << 8 & 0b0000001100000000)
	if err != nil {
		return dev.value
	}
	dev.value = int(digitalValue) - dev.midValue
	return dev.value
}

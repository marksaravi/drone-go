package mcp3008

import (
	"log"

	"periph.io/x/periph/conn/spi"
)

type mcp3008dev struct {
	spiConn       spi.Conn
	channel       int
	valueRange    byte
	digitalOffset uint16
	value         byte
}

func NewMCP3008(spiConn spi.Conn, channel int, valueRange byte, digitalOffset uint16) *mcp3008dev {
	return &mcp3008dev{
		spiConn:       spiConn,
		channel:       channel,
		valueRange:    valueRange,
		digitalOffset: digitalOffset,
	}
}

func (dev *mcp3008dev) Read() byte {
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
	digitalValue = digitalValue - dev.digitalOffset
	log.Println(digitalValue)
	return dev.value
}

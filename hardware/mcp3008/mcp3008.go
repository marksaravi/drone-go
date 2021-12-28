package mcp3008

import (
	"periph.io/x/periph/conn/spi"
)

type mcp3008dev struct {
	spiConn    spi.Conn
	channel    int
	vRef       float32
	valueRange byte
	zeroValue  float32
	value      byte
}

func NewMCP3008(spiConn spi.Conn, vRef float32, valueRange byte, channel int, zeroValue float32) *mcp3008dev {
	return &mcp3008dev{
		spiConn:    spiConn,
		channel:    channel,
		vRef:       vRef,
		valueRange: valueRange,
		zeroValue:  zeroValue,
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
	voltage := (float32(digitalValue) / 1024 * dev.vRef) - dev.zeroValue
	if voltage < 0 {
		voltage = 0
	}
	value := voltage / dev.vRef * float32(dev.valueRange)
	dev.value = byte(value)
	return dev.value
}

package mcp3008

import (
	"github.com/MarkSaravi/drone-go/connectors"
	"github.com/MarkSaravi/drone-go/modules/adcconverter"
	"periph.io/x/periph/conn/spi"
)

type MCP3008Config struct {
	SPI connectors.SPIConfig `yaml:"spi"`
}

type mcp3008dev struct {
	spiConn spi.Conn
}

func NewMCP3008(spiConn spi.Conn) adcconverter.AnalogToDigitalConverter {
	return &mcp3008dev{
		spiConn: spiConn,
	}
}

func (dev *mcp3008dev) ReadInputVoltage(channel int, vRef float32) (float32, error) {
	ch := byte(channel)
	if ch > 7 {
		ch = 0
	}
	w := []byte{0b00000001, 0b01110000 & (ch << 4), 0b00000000}
	r := []byte{0, 0, 0}
	err := dev.spiConn.Tx(w, r)
	var digitalValue uint16 = uint16(r[2]) | (uint16(r[1]) << 8 & 0b0000001100000000)
	return float32(digitalValue) / 1024 * vRef, err
}

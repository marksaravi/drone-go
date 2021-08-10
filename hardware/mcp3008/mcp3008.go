package mcp3008

import (
	"fmt"

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

func NewMCP3008(spiConn spi.Conn) adcconverter.AnalogToDigitalDevice {
	return &mcp3008dev{
		spiConn: spiConn,
	}
}

func (dev *mcp3008dev) ReadInput(channel int) int {
	//dev.spiConn.Tx()
	fmt.Println("READ INPUT")
	return 0
}

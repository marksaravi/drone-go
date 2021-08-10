package mcp3008

import (
	"fmt"

	"github.com/MarkSaravi/drone-go/modules/adcconverter"
	"periph.io/x/periph/host/sysfs"
)

type mcp3008dev struct {
	spibus *sysfs.SPI
}

func NewMCP3008(spibus *sysfs.SPI) adcconverter.AnalogToDigitalDevice {
	return &mcp3008dev{
		spibus: spibus,
	}
}

func (dev *mcp3008dev) ReadInput(channel int) int {
	fmt.Println("READ INPUT")
	return 0
}

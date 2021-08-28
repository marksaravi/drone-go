package mcp3008

import (
	"log"

	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/host/sysfs"
)

type SPIConfig struct {
	BusNumber   int `yaml:"bus_number"`
	ChipSelect  int `yaml:"chip_select"`
	Mode        int `yaml:"mode"`
	SpeedMegaHz int `yaml:"speed-mega-hz"`
}

type MCP3008Config struct {
	SPI SPIConfig `yaml:"spi"`
}

type mcp3008dev struct {
	spiConn spi.Conn
}

func NewMCP3008(busNumber int, chipSelect int, mode spi.Mode, speedMHz int) *mcp3008dev {
	spibus, _ := sysfs.NewSPI(
		busNumber,
		chipSelect,
	)
	spiConn, err := spibus.Connect(
		physic.Frequency(speedMHz)*physic.MegaHertz,
		mode,
		8,
	)
	if err != nil {
		log.Fatal(err)
	}
	return &mcp3008dev{
		spiConn: spiConn,
	}
}

func (dev *mcp3008dev) ReadInputVoltage(channel int, vRef float32, zeroValue float32) (float32, error) {
	ch := byte(channel)
	if ch > 7 {
		ch = 0
	}
	w := []byte{0b00000001, 0b10000000 | (ch << 4), 0b00000000}
	r := []byte{0, 0, 0}
	err := dev.spiConn.Tx(w, r)
	var digitalValue uint16 = uint16(r[2]) | (uint16(r[1]) << 8 & 0b0000001100000000)
	return (float32(digitalValue) / 1024 * vRef) - zeroValue, err
}

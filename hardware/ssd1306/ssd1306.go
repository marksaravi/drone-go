package ssd1306

import (
	"periph.io/x/periph/conn/i2c"
)

type Options struct {
	W             int
	H             int
	Rotated       bool
	Sequential    bool
	SwapTopBottom bool
}

type SSD1306 struct {
	I2CDev  *i2c.Dev
	Options Options
}

func (d *SSD1306) SendCommand(c []byte) error {

	return d.I2CDev.Tx(append([]byte{0x00}, c...), nil)
}

func (d *SSD1306) DisplayOff() error {

	return d.SendCommand([]byte{0xAE})
}

func (d *SSD1306) DisplayOn() error {

	return d.SendCommand([]byte{0xAF})
}

func (d *SSD1306) SendData(c []byte) error {
	return d.I2CDev.Tx(append([]byte{0x40}, c...), nil)
}

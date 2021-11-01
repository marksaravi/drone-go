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

type OLED struct {
	I2CDev  *i2c.Dev
	Options Options
}

func (d *OLED) SendCommand(c []byte) error {

	return d.I2CDev.Tx(append([]byte{0x00}, c...), nil)
}

func (d *OLED) DisplayOff() error {

	return d.SendCommand([]byte{0xAE})
}

func (d *OLED) DisplayOn() error {

	return d.SendCommand([]byte{0xAF})
}

func (d *OLED) SendData(c []byte) error {
	return d.I2CDev.Tx(append([]byte{0x40}, c...), nil)
}

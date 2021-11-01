package ssd1306

import (
	"image"

	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/devices/ssd1306/image1bit"
)

type Options struct {
	Width         int
	Height        int
	Rotated       bool
	Sequential    bool
	SwapTopBottom bool
}

var DefaultOptions = Options{
	Width:         128,
	Height:        64,
	Rotated:       false,
	Sequential:    false,
	SwapTopBottom: false,
}

type SSD1306 struct {
	I2CDev  *i2c.Dev
	Options Options
	Buffer  *image1bit.VerticalLSB
}

func NewSSD1306(dev *i2c.Dev, options Options) *SSD1306 {
	return &SSD1306{
		I2CDev:  dev,
		Options: options,
		Buffer:  image1bit.NewVerticalLSB(image.Rect(0, 0, options.Width, options.Height)),
	}
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

func (d *SSD1306) SetPixel(row, col int) {
	absIndex := (row/8)*d.Options.Width + col
	maskIndex := row % 8
	maskValue := byte(1) << byte(maskIndex)

	d.Buffer.Pix[absIndex] = d.Buffer.Pix[absIndex] | maskValue
}

func (d *SSD1306) WriteChar(charCode, x, y int) {
	var xOffset = MonoFont.width*x + 4
	var yOffset = (MonoFont.height + 5) * y
	char := MonoFont.fontData[charCode]
	for row := 0; row < MonoFont.height; row++ {
		for col := 0; col < MonoFont.width; col++ {
			if char[row][col] > 0 {
				d.SetPixel(row+yOffset, col+xOffset)
			}
		}
	}

}

func (d *SSD1306) WriteString(msg string, x, y int) {
	charCodes := []byte(msg)
	for i := 0; i < len(charCodes); i++ {
		d.WriteChar(int(charCodes[i]), x+i, y)
	}
}

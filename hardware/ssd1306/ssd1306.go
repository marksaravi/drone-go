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

func (d *SSD1306) Init() error {
	return d.SendCommand(getInitCmd(d.Options))
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

func (d *SSD1306) Draw() error {
	return d.I2CDev.Tx(append([]byte{0x40}, d.Buffer.Pix...), nil)
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

func getInitCmd(opts Options) []byte {
	// Set COM output scan direction; C0 means normal; C8 means reversed
	comScan := byte(0xC8)
	// See page 40.
	columnAddr := byte(0xA1)
	if opts.Rotated {
		// Change order both horizontally and vertically.
		comScan = 0xC0
		columnAddr = byte(0xA0)
	}
	// See page 40.
	hwLayout := byte(0x02)
	if !opts.Sequential {
		hwLayout |= 0x10
	}
	if opts.SwapTopBottom {
		hwLayout |= 0x20
	}
	// Set the max frequency. The problem with I²C is that it creates visible
	// tear down. On SPI at high speed this is not visible. Page 23 pictures how
	// to avoid tear down. For now default to max frequency.
	freq := byte(0xF0)

	// Initialize the device by fully resetting all values.
	// Page 64 has the full recommended flow.
	// Page 28 lists all the commands.
	return []byte{
		0xAE,       // Display off
		0xD3, 0x00, // Set display offset; 0
		0x40,           // Start display start line; 0
		columnAddr,     // Set segment remap; RESET is column 127.
		comScan,        //
		0xDA, hwLayout, // Set COM pins hardware configuration; see page 40
		0x81, 0xFF, // Set max contrast
		0xA4,       // Set display to use GDDRAM content
		0xA6,       // Set normal display (0xA7 for inverted 0=lit, 1=dark)
		0xD5, freq, // Set osc frequency and divide ratio; power on reset value is 0x80.
		0x8D, 0x14, // Enable charge pump regulator; page 62
		0xD9, 0xF1, // Set pre-charge period; from adafruit driver
		0xDB, 0x40, // Set Vcomh deselect level; page 32
		0x2E,                        // Deactivate scroll
		0xA8, byte(opts.Height - 1), // Set multiplex ratio (number of lines to display)
		0x20, 0x00, // Set memory addressing mode to horizontal
		0x21, 0, uint8(opts.Width - 1), // Set column address (Width)
		0x22, 0, uint8(opts.Height/8 - 1), // Set page address (Pages)
		0xAF, // Display on
	}
}

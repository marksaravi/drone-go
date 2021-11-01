package ssd1306

import (
	"fmt"
	"image"

	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
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
	options Options
	i2cdev  *i2c.Dev
	pixels  []byte
}

func NewSSD1306(address uint16, options Options) (*SSD1306, error) {
	b, err := i2creg.Open("")
	if err != nil {
		return nil, err
	}
	defer b.Close()

	dev := &i2c.Dev{Addr: address, Bus: b}
	return &SSD1306{
		i2cdev:  dev,
		options: options,
		pixels:  make([]byte, options.Width*options.Height/8),
	}, nil
}

func (oled *SSD1306) Bounds() image.Rectangle {
	return image.Rect(0, 0, oled.options.Width, oled.options.Height)
}

func (oled *SSD1306) Init() {
	oled.sendCommand(getInitCmd(oled.options))
}

func (oled *SSD1306) sendCommand(c []byte) error {

	return oled.i2cdev.Tx(append([]byte{0x00}, c...), nil)
}

func (oled *SSD1306) Draw() error {
	return oled.i2cdev.Tx(append([]byte{0x40}, oled.pixels...), nil)
}

func (oled *SSD1306) DisplayOff() error {

	return oled.sendCommand([]byte{0xAE})
}

func (oled *SSD1306) DisplayOn() error {

	return oled.sendCommand([]byte{0xAF})
}

func (oled *SSD1306) SetPixel(x, y int) {
	absIndex := (y/8)*oled.Bounds().Dx() + x
	maskIndex := y % 8
	maskValue := byte(1) << byte(maskIndex)

	oled.pixels[absIndex] = oled.pixels[absIndex] | maskValue
}

func (oled *SSD1306) ClearScreen() {
	fmt.Println(len(oled.pixels))
	for i := 0; i < len(oled.pixels); i++ {
		oled.pixels[i] = 255
	}
	oled.Draw()
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

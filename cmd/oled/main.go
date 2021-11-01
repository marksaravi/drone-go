package main

import (
	"log"

	"github.com/marksaravi/drone-go/hardware/ssd1306"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"
)

func main() {
	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Use i2creg I²C bus registry to find the first available I²C bus.
	b, err := i2creg.Open("")
	if err != nil {
		log.Fatal(err)
	}
	defer b.Close()

	d := &i2c.Dev{Addr: 0x3D, Bus: b}

	oled := ssd1306.NewSSD1306(d, ssd1306.DefaultOptions)

	err = oled.SendCommand(getInitCmd(oled.Options))
	if err != nil {
		log.Fatal(err)
	}
	oled.DisplayOn()

	for x := 0; x < 60; x++ {
		oled.SetPixel(x, x)
	}
	oled.WriteString("Hello Mark!", 0, 0)
	oled.SendData(oled.Buffer.Pix)
	// bounds := image.Rect(0, 0, oled.ops.W, oled.ops.H)
	// img := image1bit.NewVerticalLSB(bounds)

	// writeString(bounds, img, "Disconnected", 0, 1)
	// writeString(bounds, img, "Power: 54.3%", 0, 2)

}

func getInitCmd(opts ssd1306.Options) []byte {
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

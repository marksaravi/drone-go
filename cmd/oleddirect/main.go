package main

import (
	"log"

	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"
)

type options struct {
	W             int
	H             int
	Rotated       bool
	Sequential    bool
	SwapTopBottom bool
}

type OLED struct {
	conn *i2c.Dev
	ops  options
}

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

	d := &i2c.Dev{Addr: 0x3C, Bus: b}

	const W = 128
	const H = 64
	oled := OLED{
		conn: d,
		ops: options{
			W:          W,
			H:          H,
			Sequential: false,
			Rotated:    false,
		},
	}

	err = oled.sendCommand(getInitCmd(&oled.ops))
	if err != nil {
		log.Fatal(err)
	}
	buffer := make([]byte, W*H/8)
	for i := 0; i < len(buffer); i++ {
		buffer[i] = 255
	}
	err = oled.sendData(buffer)
	if err != nil {
		log.Fatal(err)
	}
	oled.displayOn()
}

func getInitCmd(opts *options) []byte {
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
		0x2E,                   // Deactivate scroll
		0xA8, byte(opts.H - 1), // Set multiplex ratio (number of lines to display)
		0x20, 0x00, // Set memory addressing mode to horizontal
		0x21, 0, uint8(opts.W - 1), // Set column address (Width)
		0x22, 0, uint8(opts.H/8 - 1), // Set page address (Pages)
		0xAF, // Display on
	}
}

func (d *OLED) sendCommand(c []byte) error {

	return d.conn.Tx(append([]byte{0x00}, c...), nil)
}

func (d *OLED) displayOff() error {

	return d.sendCommand([]byte{0xAE})
}
func (d *OLED) displayOn() error {

	return d.sendCommand([]byte{0xAF})
}

func (d *OLED) sendData(c []byte) error {
	return d.conn.Tx(append([]byte{0x40}, c...), nil)
}

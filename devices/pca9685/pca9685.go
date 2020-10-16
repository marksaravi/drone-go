package pca9685

import (
	"time"

	"github.com/MarkSaravi/drone-go/drivers/i2c"
)

const PCA9685Address = 0x40

const (
	PCA9685_MODE1        = 0x00
	PCA9685_MODE2        = 0x01
	PCA9685_PRESCALE     = 0xFE
	PCA9685_SUBADR1      = 0x02
	PCA9685_SUBADR2      = 0x03
	PCA9685_SUBADR3      = 0x04
	PCA9685_LED0_ON_L    = 0x06
	PCA9685_LED0_ON_H    = 0x07
	PCA9685_LED0_OFF_L   = 0x08
	PCA9685_LED0_OFF_H   = 0x09
	PCA9685_ALLLED_ON_L  = 0xFA
	PCA9685_ALLLED_ON_H  = 0xFB
	PCA9685_ALLLED_OFF_L = 0xFC
	PCA9685_ALLLED_OFF_H = 0xFD

	PCA9685_RESTART = 0x80
	PCA9685_SLEEP   = 0x10
	PCA9685_ALLCALL = 0x01
	PCA9685_INVRT   = 0x10
	PCA9685_OUTDRV  = 0x04
)

// PCA9685Driver is struct for PCA9685
type PCA9685Driver struct {
	name       string
	address    uint8
	connection *i2c.Connection
	frequency  float32
}

// NewPCA9685Driver creates new PCA9685 driver
func NewPCA9685Driver(address uint8, connection *i2c.Connection) (*PCA9685Driver, error) {
	return &PCA9685Driver{
		name:       "PCA9685",
		address:    address,
		connection: connection,
	}, nil
}

func (d *PCA9685Driver) readByte(offset uint8) (b uint8, err error) {
	return d.connection.ReadByte(d.address, offset)
}

func (d *PCA9685Driver) writeByte(offset uint8, b uint8) (err error) {
	return d.connection.WriteByte(d.address, offset, b)
}

func (d *PCA9685Driver) writeAddress(offset uint8) (err error) {
	return d.connection.WriteBytes(d.address, offset)
}

// Close closes the i2c connection
func (d *PCA9685Driver) Close() {
	d.connection.Close()
}

// Halt stops the device
func (d *PCA9685Driver) Halt() (err error) {
	err = d.writeByte(PCA9685_ALLLED_OFF_H, 0x10)
	return
}

// setPWM sets pwm for a channel
func (d *PCA9685Driver) setPWM(channel int, on uint16, off uint16) (err error) {
	if err := d.writeByte(byte(PCA9685_LED0_ON_L+4*channel), byte(on)&0xFF); err != nil {
		return err
	}

	if err := d.writeByte(byte(PCA9685_LED0_ON_H+4*channel), byte(on>>8)); err != nil {
		return err
	}

	if err := d.writeByte(byte(PCA9685_LED0_OFF_L+4*channel), byte(off)&0xFF); err != nil {
		return err
	}

	if err := d.writeByte(byte(PCA9685_LED0_OFF_H+4*channel), byte(off>>8)); err != nil {
		return err
	}

	return
}

// SetPWMFreq sets the PWM frequency in Hz
func (d *PCA9685Driver) setPWMFreq(freq float32) error {
	d.frequency = freq
	// IC oscillator frequency is 25 MHz
	var prescalevel float32 = 25000000
	// Find frequency of PWM waveform
	prescalevel /= 4096
	// Ratio between desired frequency and maximum
	prescalevel /= freq
	prescalevel--
	// Round value to nearest whole
	prescale := byte(prescalevel + 0.5)

	if err := d.writeAddress(byte(PCA9685_MODE1)); err != nil {
		return err
	}
	oldmode, err := d.readByte(byte(PCA9685_MODE1))
	if err != nil {
		return err
	}

	// Put oscillator in sleep mode, clear bit 7 here to avoid overwriting
	// previous setting
	newmode := (oldmode & 0x7F) | 0x10
	if err := d.writeByte(byte(PCA9685_MODE1), byte(newmode)); err != nil {
		return err
	}
	// Write prescaler value
	if err := d.writeByte(byte(PCA9685_PRESCALE), prescale); err != nil {
		return err
	}
	// Put back to old settings
	if err := d.writeByte(byte(PCA9685_MODE1), byte(oldmode)); err != nil {
		return err
	}

	time.Sleep(5 * time.Millisecond)

	// Enable response to All Call address, enable auto-increment, clear restart
	if err := d.writeByte(byte(PCA9685_MODE1), byte(oldmode|0x80)); err != nil {
		return err
	}

	return nil
}

func (d *PCA9685Driver) setAllPWM(on uint16, off uint16) (err error) {
	if err := d.writeByte(byte(PCA9685_ALLLED_ON_L), byte(on)&0xFF); err != nil {
		return err
	}

	if err := d.writeByte(byte(PCA9685_ALLLED_ON_H), byte(on>>8)); err != nil {
		return err
	}

	if err := d.writeByte(byte(PCA9685_ALLLED_OFF_L), byte(off)&0xFF); err != nil {
		return err
	}

	if err := d.writeByte(byte(PCA9685_ALLLED_OFF_H), byte(off>>8)); err != nil {
		return err
	}

	return
}

func (d *PCA9685Driver) Start(frequency float32) error {
	if err := d.setAllPWM(0, 0); err != nil {
		return err
	}

	if err := d.writeByte(PCA9685_MODE2, PCA9685_OUTDRV); err != nil {
		return err
	}

	if err := d.writeByte(PCA9685_MODE1, PCA9685_ALLCALL); err != nil {
		return err
	}

	time.Sleep(5 * time.Millisecond)

	if err := d.writeAddress(byte(PCA9685_MODE1)); err != nil {
		return err
	}
	oldmode, err := d.readByte(PCA9685_MODE1)
	if err != nil {
		return err
	}
	oldmode = oldmode &^ byte(PCA9685_SLEEP)

	if err := d.writeByte(PCA9685_MODE1, oldmode); err != nil {
		return err
	}

	time.Sleep(5 * time.Millisecond)
	d.setPWMFreq(frequency)
	return err
}

func (d *PCA9685Driver) SetPulseWidth(channel int, pulseWidth float32) {
	period := float32(1) / d.frequency
	on := pulseWidth / period * 4096
	d.setPWM(channel, 0, uint16(on))
}

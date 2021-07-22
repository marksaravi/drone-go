package pca9685

import (
	"time"

	"github.com/MarkSaravi/drone-go/connectors/i2c"
)

//PCA9685Address is i2c address of device
const PCA9685Address = 0x40

//Addresses
const (
	PCA9685Mode1      = 0x00
	PCA9685Mode2      = 0x01
	PCA9685Prescale   = 0xFE
	PCA9685Subaddr1   = 0x02
	PCA9685Subaddr2   = 0x03
	PCA9685Subaddr3   = 0x04
	PCA9685LED0OnL    = 0x06
	PCA9685LED0OnH    = 0x07
	PCA9685LED0OffL   = 0x08
	PCA9685LED0OffH   = 0x09
	PCA9685AlliedOnL  = 0xFA
	PCA9685AlliedOnH  = 0xFB
	PCA9685AlliedOffL = 0xFC
	PCA9685AlliedOffH = 0xFD

	PCA9685Restart = 0x80
	PCA9685Sleep   = 0x10
	PCA9685AllCall = 0x01
	PCA9685Invert  = 0x10
	PCA9685OutDrv  = 0x04
)

const (
	Frequency       float32 = 400
	MinPW           float32 = 0.001
	MaxPW           float32 = 0.002
	MaxApplicablePW float32 = MinPW + (MaxPW-MinPW)*0.1
)

//PCA9685 is struct for PCA9685
type PCA9685 struct {
	name       string
	address    uint8
	connection *i2c.Connection
	frequency  float32
}

func (d *PCA9685) readByte(offset uint8) (b uint8, err error) {
	return d.connection.ReadByte(d.address, offset)
}

func (d *PCA9685) writeByte(offset uint8, b uint8) (err error) {
	return d.connection.WriteByte(d.address, offset, b)
}

func (d *PCA9685) writeAddress(offset uint8) (err error) {
	return d.connection.WriteBytes(d.address, offset)
}

func getOffTime(frequency float32, pulseWidth float32) (on uint16, off uint16) {
	period := float32(1) / frequency
	on = 0
	off = uint16(pulseWidth / period * 4096)
	return
}

func (d *PCA9685) setPWM(channel int, pulseWidth float32) (err error) {
	on, off := getOffTime(d.frequency, pulseWidth)

	if err := d.writeByte(byte(PCA9685LED0OnL+4*channel), byte(on)&0xFF); err != nil {
		return err
	}

	if err := d.writeByte(byte(PCA9685LED0OnH+4*channel), byte(on>>8)); err != nil {
		return err
	}

	if err := d.writeByte(byte(PCA9685LED0OffL+4*channel), byte(off)&0xFF); err != nil {
		return err
	}

	if err := d.writeByte(byte(PCA9685LED0OffH+4*channel), byte(off>>8)); err != nil {
		return err
	}

	return
}

// SetPWMFreq sets the PWM frequency in Hz
func (d *PCA9685) setFrequency(freq float32) error {
	d.frequency = freq
	// IC oscillator frequency is 25 MHz
	var prescalevel float32 = 24576000
	// Find frequency of PWM waveform
	prescalevel /= 4096
	// Ratio between desired frequency and maximum
	prescalevel /= freq
	// prescalevel
	// Round value to nearest whole
	prescale := byte(prescalevel)

	if err := d.writeAddress(byte(PCA9685Mode1)); err != nil {
		return err
	}
	oldmode, err := d.readByte(byte(PCA9685Mode1))
	if err != nil {
		return err
	}

	// Put oscillator in sleep mode, clear bit 7 here to avoid overwriting
	// previous setting
	newmode := (oldmode & 0x7F) | 0x10
	if err := d.writeByte(byte(PCA9685Mode1), byte(newmode)); err != nil {
		return err
	}
	// Write prescaler value
	if err := d.writeByte(byte(PCA9685Prescale), prescale); err != nil {
		return err
	}
	// Put back to old settings
	if err := d.writeByte(byte(PCA9685Mode1), byte(oldmode)); err != nil {
		return err
	}

	time.Sleep(5 * time.Millisecond)

	// Enable response to All Call address, enable auto-increment, clear restart
	if err := d.writeByte(byte(PCA9685Mode1), byte(oldmode|0x80)); err != nil {
		return err
	}

	return nil
}

func (d *PCA9685) setAllPWM(pulseWidth float32) (err error) {
	on, off := getOffTime(d.frequency, pulseWidth)
	if err := d.writeByte(byte(PCA9685AlliedOnL), byte(on)&0xFF); err != nil {
		return err
	}

	if err := d.writeByte(byte(PCA9685AlliedOnH), byte(on>>8)); err != nil {
		return err
	}

	if err := d.writeByte(byte(PCA9685AlliedOffL), byte(off)&0xFF); err != nil {
		return err
	}

	if err := d.writeByte(byte(PCA9685AlliedOffH), byte(off>>8)); err != nil {
		return err
	}

	return
}

//Start starts the device with a frequency
func (d *PCA9685) Start() error {
	if err := d.setAllPWM(0); err != nil {
		return err
	}

	if err := d.writeByte(PCA9685Mode2, PCA9685OutDrv); err != nil {
		return err
	}

	if err := d.writeByte(PCA9685Mode1, PCA9685AllCall); err != nil {
		return err
	}

	time.Sleep(5 * time.Millisecond)

	if err := d.writeAddress(byte(PCA9685Mode1)); err != nil {
		return err
	}
	oldmode, err := d.readByte(PCA9685Mode1)
	if err != nil {
		return err
	}
	oldmode = oldmode &^ byte(PCA9685Sleep)

	if err := d.writeByte(PCA9685Mode1, oldmode); err != nil {
		return err
	}

	time.Sleep(5 * time.Millisecond)
	d.setFrequency(Frequency)
	return err
}

//SetPulseWidth sets PWM for a channel
func (d *PCA9685) SetThrottle(motor int, throttle float32) {
	d.SetPulseWidth(motor, MinPW+throttle*(MaxPW-MinPW))
}

//SetPulseWidth sets PWM for a channel
func (d *PCA9685) SetPulseWidth(channel int, pulseWidth float32) {
	d.setPWM(channel, pulseWidth)
}

// Close closes the i2c connection
func (d *PCA9685) Close() {
	d.connection.Close()
}

// Halt stops the device
func (d *PCA9685) Halt() (err error) {
	err = d.writeByte(PCA9685AlliedOffH, 0x10)
	return
}

//StopAll stops all channels
func (d *PCA9685) StopAll() {
	d.setAllPWM(0)
}

// NewPCA9685Driver creates new PCA9685 driver
func NewPCA9685Driver(address uint8, connection *i2c.Connection) (*PCA9685, error) {
	return &PCA9685{
		name:       "PCA9685",
		address:    address,
		connection: connection,
	}, nil
}

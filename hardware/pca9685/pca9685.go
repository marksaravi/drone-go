package pca9685

import (
	"fmt"
	"time"

	"github.com/marksaravi/drone-go/models"
	"github.com/marksaravi/drone-go/utils"
	"periph.io/x/periph/conn/i2c"
)

type powerbreaker interface {
	Connect()
	Disconnect()
}

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
	Frequency            float64 = 384
	MinPW                float64 = 0.000995
	MaxPW                float64 = 0.00199
	MaxAllowedThrottle   float64 = 65
	MinSafeStartThrottle float64 = 5
)

type pca9685Dev struct {
	name               string
	address            uint8
	connection         *i2c.Dev
	frequency          float64
	channelMappings    map[int]int
	maxThrottle        float64
	controlVariableMin float64
	controlVariableMax float64
	safeStartThrottle  float64
	throttle           float64
}

type PCA9685Settings struct {
	Connection           *i2c.Dev
	ChannelMappings      map[int]int
	SafeStartThrottle    float64
	MaxThrottle          float64
	ControlVariableRange float64
}

// NewPCA9685Driver creates new pca9685Dev driver
func NewPCA9685(settings PCA9685Settings) (*pca9685Dev, error) {
	validateSettings(&settings)
	dev := &pca9685Dev{
		name:               "pca9685Dev",
		address:            PCA9685Address,
		connection:         settings.Connection,
		safeStartThrottle:  settings.SafeStartThrottle,
		maxThrottle:        settings.MaxThrottle,
		channelMappings:    settings.ChannelMappings,
		controlVariableMin: -settings.ControlVariableRange / 2,
		controlVariableMax: settings.ControlVariableRange / 2,
		throttle:           0,
	}
	dev.init()
	return dev, nil
}

func (d *pca9685Dev) SetThrottles(throttles models.Throttles) {
	for i := 0; i < len(throttles.ControlVariables); i++ {
		channel := d.channelMappings[i]
		d.SetThrottle(channel, throttles.Throttle, throttles.ControlVariables[i])
	}
}

func (d *pca9685Dev) SetThrottle(channel int, throttle, controlVariable float64) {
	t := d.validateThrottle(throttle, controlVariable)
	pulseWidth := MinPW + t/100*(MaxPW-MinPW)
	d.setPWM(channel, pulseWidth)
}

func (d *pca9685Dev) validateThrottle(throttle, controlVariable float64) float64 {
	t := utils.ApplyLimit(throttle, 0, d.maxThrottle, true)
	v := utils.ApplyLimit(controlVariable, d.controlVariableMin, d.controlVariableMax, true)
	if t <= d.safeStartThrottle {
		v = 0
	}

	return utils.ApplyLimit(v+t, 0, MaxAllowedThrottle, false)
}

//Calibrate
func Calibrate(i2cConn *i2c.Dev, powerbreaker powerbreaker, mappings map[int]int) {
	pwmDev, err := NewPCA9685(PCA9685Settings{Connection: i2cConn, SafeStartThrottle: MinSafeStartThrottle, MaxThrottle: MinSafeStartThrottle + 1, ChannelMappings: mappings})
	if err != nil {
		fmt.Println(err)
		return
	}

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("setting max pulse width: ", MaxPW)
	fmt.Println("turn on ESCs")
	pwmDev.setAllPWM(MaxPW)
	time.Sleep(1 * time.Second)
	powerbreaker.Connect()
	time.Sleep(12 * time.Second)
	fmt.Println("setting min pulse width: ", MinPW)
	pwmDev.setAllPWM(MinPW)
	time.Sleep(12 * time.Second)
	fmt.Println("turn off ESCs")
	powerbreaker.Disconnect()
	time.Sleep(1 * time.Second)
	pwmDev.setAllPWM(0)
}

func validateSettings(settings *PCA9685Settings) {
	if settings.SafeStartThrottle < MinSafeStartThrottle {
		panic(fmt.Errorf("safe start throttle must be a positive number and more than %6.1f", MinSafeStartThrottle))
	}
	if settings.MaxThrottle <= settings.SafeStartThrottle {
		panic(fmt.Errorf("max throttle must be a more than %6.1f", settings.SafeStartThrottle))
	}
	if settings.MaxThrottle+settings.ControlVariableRange/2 > MaxAllowedThrottle {
		panic(fmt.Errorf("throttle + control variable must be less than %6.1f", MaxAllowedThrottle))
	}
}

func (d *pca9685Dev) readByte(offset uint8) (uint8, error) {
	read := make([]byte, 1)
	write := []byte{offset}
	err := d.connection.Tx(write, read)
	return read[0], err
}

func (d *pca9685Dev) writeByte(offset uint8, b uint8) error {
	write := []byte{offset, b}
	_, err := d.connection.Write(write)
	return err
}

func getOffTime(frequency float64, pulseWidth float64) (on uint16, off uint16) {
	period := float64(1) / frequency
	on = 0
	off = uint16(pulseWidth / period * 4096)
	return
}

func (d *pca9685Dev) setPWM(channel int, pulseWidth float64) (err error) {
	on, off := getOffTime(d.frequency, pulseWidth)
	addresses := []byte{
		byte(PCA9685LED0OnL + 4*channel),
		byte(PCA9685LED0OnH + 4*channel),
		byte(PCA9685LED0OffL + 4*channel),
		byte(PCA9685LED0OffH + 4*channel),
	}
	values := []byte{
		byte(on) & 0xFF,
		byte(on >> 8),
		byte(off) & 0xFF,
		byte(off >> 8),
	}

	for i := 0; i < 4; i++ {
		w := []byte{addresses[i], values[i]}
		if err := d.connection.Tx(w, nil); err != nil {
			return err
		}
	}
	return
}

// SetPWMFreq sets the PWM frequency in Hz
func (d *pca9685Dev) setFrequency(freq float64) error {
	d.frequency = freq
	// IC oscillator frequency is 25 MHz
	var prescalevel float64 = 24576000
	// Find frequency of PWM waveform
	prescalevel /= 4096
	// Ratio between desired frequency and maximum
	prescalevel /= freq
	// prescalevel
	// Round value to nearest whole
	prescale := byte(prescalevel)

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

func (d *pca9685Dev) setAllPWM(pulseWidth float64) (err error) {
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
func (d *pca9685Dev) init() error {
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

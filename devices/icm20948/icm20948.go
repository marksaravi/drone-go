package icm20948

import (
	"time"

	"github.com/MarkSaravi/drone-go/utils"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/host"
	"periph.io/x/periph/host/sysfs"
)

const (
	BANK0 uint16 = 0 << 8
	BANK1 uint16 = 1 << 8
	BANK2 uint16 = 2 << 8
	BANK3 uint16 = 3 << 8
)

const (
	REG_BANK_SEL byte = 0x7F

	// BANK0
	WHO_AM_I     uint16 = BANK0 | 0x0
	LP_CONFIG    uint16 = BANK0 | 0x5
	PWR_MGMT_1   uint16 = BANK0 | 0x6
	PWR_MGMT_2   uint16 = BANK0 | 0x7
	INT_ENABLE_3 uint16 = BANK0 | 0x13
	ACCEL_XOUT_H uint16 = BANK0 | 0x2D
	ACCEL_XOUT_L uint16 = BANK0 | 0x2E
	ACCEL_YOUT_H uint16 = BANK0 | 0x2F
	ACCEL_YOUT_L uint16 = BANK0 | 0x30
	ACCEL_ZOUT_H uint16 = BANK0 | 0x31
	ACCEL_ZOUT_L uint16 = BANK0 | 0x32
	GYRO_XOUT_H  uint16 = BANK0 | 0x33
	GYRO_XOUT_L  uint16 = BANK0 | 0x34
	GYRO_YOUT_H  uint16 = BANK0 | 0x35
	GYRO_YOUT_L  uint16 = BANK0 | 0x36
	GYRO_ZOUT_H  uint16 = BANK0 | 0x37
	GYRO_ZOUT_L  uint16 = BANK0 | 0x38

	// BANK1
	XA_OFFS_H uint16 = BANK1 | 0x14

	// BANK2
	GYRO_SMPLRT_DIV uint16 = BANK2 | 0x0
	GYRO_CONFIG_1   uint16 = BANK2 | 0x1
	GYRO_CONFIG_2   uint16 = BANK2 | 0x2
	ZG_OFFS_USRL    uint16 = BANK2 | 0x8
	ACCEL_CONFIG    uint16 = BANK2 | 0x14
	ACCEL_CONFIG_2  uint16 = BANK2 | 0x15
	MOD_CTRL_USR    uint16 = BANK2 | 0x54
)

func reg(reg uint16) *Register {
	return &Register{
		address: byte(reg),
		bank:    byte(reg >> 8),
	}
}

// Register is the address and bank of the register
type Register struct {
	address byte
	bank    byte
}

// Device is icm20948 mems
type Device struct {
	*sysfs.SPI
	spi.Conn
	regbank                  byte
	accelerometerSensitivity int
}

var accelerometerSensitivity = make(map[int]float64)

func init() {
	accelerometerSensitivity[0] = 16384
	accelerometerSensitivity[1] = 8192
	accelerometerSensitivity[2] = 4096
	accelerometerSensitivity[3] = 2048
	host.Init()
}

// NewRaspberryPiICM20948Driver creates ICM20948 driver for raspberry pi
func NewRaspberryPiICM20948Driver(busNumber int, chipSelect int) (*Device, error) {
	d, err := sysfs.NewSPI(busNumber, chipSelect)
	if err != nil {
		return nil, err
	}
	conn, err := d.Connect(7*physic.MegaHertz, spi.Mode3, 8)
	if err != nil {
		return nil, err
	}
	dev := Device{
		SPI:     d,
		Conn:    conn,
		regbank: 0xFF,
	}
	return &dev, nil
}

func (dev *Device) readReg(address byte, len int) ([]byte, error) {
	w := make([]byte, len+1)
	r := make([]byte, len+1)
	w[0] = (address & 0x7F) | 0x80
	err := dev.Conn.Tx(w, r)
	return r[1:], err
}

func (dev *Device) writeReg(address byte, data ...byte) error {
	if len(data) == 0 {
		return nil
	}
	w := append([]byte{address & 0x7F}, data...)
	err := dev.Conn.Tx(w, nil)
	return err
}

func (dev *Device) selRegisterBank(regbank byte) error {
	if regbank == dev.regbank {
		return nil
	}
	dev.regbank = regbank
	return dev.writeReg(REG_BANK_SEL, (regbank<<4)&0x30)
}

func (dev *Device) readRegister(register uint16, len int) ([]byte, error) {
	reg := reg(register)
	dev.selRegisterBank(reg.bank)
	return dev.readReg(reg.address, len)
}

func (dev *Device) writeRegister(register uint16, data ...byte) error {
	if len(data) == 0 {
		return nil
	}
	reg := reg(register)
	dev.selRegisterBank(reg.bank)
	return dev.writeReg(reg.address, data...)
}

// WhoAmI return value for ICM-20948 is 0xEA
func (dev *Device) WhoAmI() (name string, id byte, err error) {
	name = "ICM-20948"
	data, err := dev.readRegister(WHO_AM_I, 1)
	id = data[0]
	return
}

// GetDeviceConfig reads device configurations
func (dev *Device) GetDeviceConfig() ([]byte, error) {
	data, err := dev.readRegister(LP_CONFIG, 3)
	return data, err
}

// SetDeviceConfig applies initial configurations for device
func (dev *Device) SetDeviceConfig() error {
	// Reset settings to default
	err := dev.writeRegister(PWR_MGMT_1, 0b10000000)
	time.Sleep(50 * time.Millisecond) // wait for taking effect
	data, err := dev.readRegister(PWR_MGMT_1, 1)
	const nosleep byte = 0b10111111
	config := byte(data[0] & nosleep)
	const accGyro byte = 0b00000000
	err = dev.writeRegister(PWR_MGMT_1, config, accGyro)
	time.Sleep(50 * time.Millisecond) // wait for taking effect
	return err
}

// ReadRawData reads all Accl and Gyro data
func (dev *Device) ReadRawData() ([]byte, error) {
	return dev.readRegister(ACCEL_XOUT_H, 12)
}

// ReadData reads Accelerometer and Gyro data
func (dev *Device) ReadData() (accX, accY, accZ, gyroX, gyroY, gyroZ float64, err error) {
	data, err := dev.ReadRawData()
	accSens := accelerometerSensitivity[dev.accelerometerSensitivity]
	accX = float64(utils.TowsComplementBytesToInt(data[0], data[1])) / accSens
	accY = float64(utils.TowsComplementBytesToInt(data[2], data[3])) / accSens
	accZ = float64(utils.TowsComplementBytesToInt(data[4], data[5])) / accSens
	gyroX = float64(utils.TowsComplementBytesToInt(data[6], data[7])) / accSens
	gyroY = float64(utils.TowsComplementBytesToInt(data[8], data[9])) / accSens
	gyroZ = float64(utils.TowsComplementBytesToInt(data[10], data[11])) / accSens
	return
}

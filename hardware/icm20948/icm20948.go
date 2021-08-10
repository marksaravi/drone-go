package icm20948

import (
	"time"

	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/host/sysfs"
)

// XYZ is X, Y, Z data
type XYZ struct {
	X, Y, Z float64
}

type Offsets struct {
	X int16 `yaml:"X"`
	Y int16 `yaml:"Y"`
	Z int16 `yaml:"Z"`
}

// Config is the generic configuration
type Config interface {
}

// Sensor is devices that read data in x, y, z format
type Sensor struct {
	Type   string
	Config Config
}

// Register is the address and bank of the Register
type Register struct {
	Address uint8
	Bank    uint8
}

type SensorData struct {
	Error error
	Data  XYZ
}

type ICM20948Config struct {
	BusNumber  int                 `yaml:"bus_number"`
	ChipSelect int                 `yaml:"chip_select"`
	AccConfig  AccelerometerConfig `yaml:"accelerometer"`
	GyroConfig GyroscopeConfig     `yaml:"gyroscope"`
}

// memsICM20948 is icm20948 mems
type memsICM20948 struct {
	Name string
	*sysfs.SPI
	spi.Conn
	regbank uint8
	acc     Sensor
	gyro    Sensor
	mag     Sensor
}

func reg(reg uint16) *Register {
	return &Register{
		Address: uint8(reg),
		Bank:    uint8(reg >> 8),
	}
}

var accelerometerSensitivity = make(map[string]float64)
var gyroFullScale = make(map[string]float64)

func init() {
	accelerometerSensitivity["2g"] = SENSITIVITY_2G
	accelerometerSensitivity["4g"] = SENSITIVITY_4G
	accelerometerSensitivity["8g"] = SENSITIVITY_8G
	accelerometerSensitivity["16g"] = SENSITIVITY_16G

	gyroFullScale["250dps"] = GYRO_SCALE_250DPS
	gyroFullScale["500dps"] = GYRO_SCALE_500DPS
	gyroFullScale["1000dps"] = GYRO_SCALE_1000DPS
	gyroFullScale["2000dps"] = GYRO_SCALE_2000DPS
}

// NewICM20948Driver creates ICM20948 driver for raspberry pi
func NewICM20948Driver(config ICM20948Config) (*memsICM20948, error) {
	d, err := sysfs.NewSPI(config.BusNumber, config.ChipSelect)
	if err != nil {
		return nil, err
	}
	conn, err := d.Connect(7*physic.MegaHertz, spi.Mode3, 8)
	if err != nil {
		return nil, err
	}

	dev := memsICM20948{
		Name:    "ICM20948",
		SPI:     d,
		Conn:    conn,
		regbank: 0xFF,
		acc: Sensor{
			Type:   ACCELEROMETER,
			Config: config.AccConfig,
		},
		gyro: Sensor{
			Type:   GYROSCOPE,
			Config: config.GyroConfig,
		},
		mag: Sensor{
			Type: MAGNETOMETER,
		},
	}
	dev.initDevice()
	return &dev, nil
}

func (dev *memsICM20948) readReg(address uint8, len int) ([]uint8, error) {
	w := make([]uint8, len+1)
	r := make([]uint8, len+1)
	w[0] = (address & 0x7F) | 0x80
	err := dev.Conn.Tx(w, r)
	return r[1:], err
}

func (dev *memsICM20948) writeReg(address uint8, data ...uint8) error {
	if len(data) == 0 {
		return nil
	}
	w := append([]uint8{address & 0x7F}, data...)
	err := dev.Conn.Tx(w, nil)
	return err
}

func (dev *memsICM20948) selRegisterBank(regbank uint8) error {
	if regbank == dev.regbank {
		return nil
	}
	dev.regbank = regbank
	return dev.writeReg(REG_BANK_SEL, regbank<<4)
}

func (dev *memsICM20948) readRegister(register uint16, len int) ([]uint8, error) {
	reg := reg(register)
	dev.selRegisterBank(reg.Bank)
	return dev.readReg(reg.Address, len)
}

func (dev *memsICM20948) writeRegister(register uint16, data ...uint8) error {
	if len(data) == 0 {
		return nil
	}
	reg := reg(register)
	dev.selRegisterBank(reg.Bank)
	return dev.writeReg(reg.Address, data...)
}

func (dev *memsICM20948) initDevice() error {
	// Reset settings to default
	err := dev.writeRegister(PWR_MGMT_1, 0b10000000)
	if err != nil {
		return err
	}
	time.Sleep(50 * time.Millisecond) // wait for taking effect
	if err != nil {
		return err
	}
	// No low power mode, enabling everything with 20Mhz clock
	err = dev.writeRegister(PWR_MGMT_1, 0b00000001, 0b00000000)
	if err != nil {
		return err
	}
	time.Sleep(50 * time.Millisecond) // wait for starting
	err = dev.InitAccelerometer()
	if err != nil {
		return err
	}
	time.Sleep(50 * time.Millisecond) // wait for starting
	err = dev.InitGyroscope()
	time.Sleep(50 * time.Millisecond) // wait for starting
	return err
}

// readSensorsRawData reads all Accl and Gyro data
func (dev *memsICM20948) readSensorsRawData() ([]uint8, error) {
	return dev.readRegister(ACCEL_XOUT_H, 12)
}

// ReadSensors reads Accelerometer and Gyro data
func (dev *memsICM20948) ReadSensors() (
	acc SensorData,
	gyro SensorData,
	mag SensorData,
	err error) {
	data, err := dev.readSensorsRawData()

	if err != nil {
		return
	}
	accData, accErr := dev.processAccelerometerData(data)
	gyroData, gyroErr := dev.processGyroscopeData(data[6:])

	acc = SensorData{
		Error: accErr,
		Data:  accData,
	}
	gyro = SensorData{
		Error: gyroErr,
		Data:  gyroData,
	}
	mag = SensorData{
		Error: nil,
		Data:  XYZ{X: 0, Y: 0, Z: 0},
	}
	return
}

// towsComplementUint8ToInt16 converts 2's complement H and L uint8 to signed int16
func towsComplementUint8ToInt16(h, l uint8) int16 {
	var h16 uint16 = uint16(h)
	var l16 uint16 = uint16(l)

	return int16((h16 << 8) | l16)
}

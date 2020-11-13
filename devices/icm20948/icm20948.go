package icm20948

import (
	"time"

	"github.com/MarkSaravi/drone-go/modules/mpu/threeaxissensore"
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

type threeAxis struct {
	data     threeaxissensore.Data
	prevData threeaxissensore.Data
	dataDiff float64
	config   threeaxissensore.Config
}

// DeviceConfig is the configuration for the device
type DeviceConfig struct {
}

// AccelerometerConfig is the configurations for Accelerometer
type AccelerometerConfig struct {
	Sensitivity int
}

// GyroscopeConfig is the configuration for Gyroscope
type GyroscopeConfig struct {
	FullScale int
}

// Device is icm20948 mems
type Device struct {
	*sysfs.SPI
	spi.Conn
	regbank     byte
	lastReading int64
	duration    int64
	config      DeviceConfig
	acc         threeAxis
	gyro        threeAxis
}

var accelerometerSensitivity = make(map[int]float64)
var gyroFullScale = make(map[int]float64)

func init() {
	accelerometerSensitivity[0] = 16384
	accelerometerSensitivity[1] = 8192
	accelerometerSensitivity[2] = 4096
	accelerometerSensitivity[3] = 2048

	gyroFullScale[0] = 250
	gyroFullScale[1] = 500
	gyroFullScale[2] = 1000
	gyroFullScale[3] = 2000

	host.Init()
}

// NewRaspberryPiICM20948Driver creates ICM20948 driver for raspberry pi
func NewRaspberryPiICM20948Driver(
	busNumber int,
	chipSelect int,
	config DeviceConfig,
	accConfig AccelerometerConfig,
	gyroConfig GyroscopeConfig,

) (*Device, error) {
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
		acc: threeAxis{
			data:     threeaxissensore.Data{X: 0, Y: 0, Z: 0},
			prevData: threeaxissensore.Data{X: 0, Y: 0, Z: 0},
			dataDiff: 0,
			config:   accConfig,
		},
		gyro: threeAxis{
			data:     threeaxissensore.Data{X: 0, Y: 0, Z: 0},
			prevData: threeaxissensore.Data{X: 0, Y: 0, Z: 0},
			dataDiff: 0,
			config:   gyroConfig,
		},
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
func (dev *Device) GetDeviceConfig() (
	config threeaxissensore.Config,
	accConfig threeaxissensore.Config,
	gyroConfig threeaxissensore.Config,
	err error) {
	// data, err := dev.readRegister(LP_CONFIG, 3)
	config = DeviceConfig{}
	accConfig, err = dev.getAccConfig()
	gyroConfig, err = dev.getGyroConfig()
	return
}

// InitDevice applies initial configurations to device
func (dev *Device) InitDevice() error {
	// Reset settings to default
	err := dev.writeRegister(PWR_MGMT_1, 0b10000000)
	time.Sleep(50 * time.Millisecond) // wait for taking effect
	data, err := dev.readRegister(PWR_MGMT_1, 1)
	const nosleep byte = 0b10111111
	config := byte(data[0] & nosleep)
	const accGyro byte = 0b00000000
	err = dev.writeRegister(PWR_MGMT_1, config, accGyro)
	time.Sleep(50 * time.Millisecond) // wait for taking effect
	err = dev.InitAccelerometer()
	time.Sleep(50 * time.Millisecond) // wait for taking effect
	err = dev.InitGyroscope()
	time.Sleep(50 * time.Millisecond) // wait for taking effect
	return err
}

// ReadRawData reads all Accl and Gyro data
func (dev *Device) ReadRawData() ([]byte, error) {
	return dev.readRegister(ACCEL_XOUT_H, 12)
}

// Start starts device
func (dev *Device) Start() {
	dev.lastReading = time.Now().UnixNano()
}

// ReadData reads Accelerometer and Gyro data
func (dev *Device) ReadData() (acc threeaxissensore.Data, gyro threeaxissensore.Data, err error) {
	data, err := dev.ReadRawData()
	now := time.Now().UnixNano()
	dev.duration = dev.lastReading - now
	dev.lastReading = now
	dev.processAccelerometerData(data)
	dev.processGyroscopeData(data[6:])
	return dev.GetAcc().GetData(), dev.GetGyro().GetData(), err
}

func (a *threeAxis) GetConfig() threeaxissensore.Config {
	return a.config
}

func (a *threeAxis) SetConfig(config threeaxissensore.Config) {
	a.config = config
}

func (a *threeAxis) GetData() threeaxissensore.Data {
	return a.data
}

func (a *threeAxis) SetData(x, y, z float64) {
	a.prevData = a.data
	a.data = threeaxissensore.Data{
		X: x,
		Y: y,
		Z: z,
	}
	a.dataDiff = utils.CalcVectorLen(a.data) - utils.CalcVectorLen(a.prevData)
}

func (a *threeAxis) GetDiff() float64 {
	return a.dataDiff
}

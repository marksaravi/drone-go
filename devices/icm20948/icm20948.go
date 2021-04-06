package icm20948

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/modules/imu"
	"github.com/MarkSaravi/drone-go/types"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/host"
	"periph.io/x/periph/host/sysfs"
)

func reg(reg uint16) *Register {
	return &Register{
		address: uint8(reg),
		bank:    uint8(reg >> 8),
	}
}

var accelerometerSensitivity = make(map[int]float64)
var gyroFullScale = make(map[int]float64)

func init() {
	accelerometerSensitivity[0] = SENSITIVITY_0
	accelerometerSensitivity[1] = SENSITIVITY_1
	accelerometerSensitivity[2] = SENSITIVITY_2
	accelerometerSensitivity[3] = SENSITIVITY_3

	gyroFullScale[0] = SCALE_0
	gyroFullScale[1] = SCALE_1
	gyroFullScale[2] = SCALE_2
	gyroFullScale[3] = SCALE_3

	host.Init()
}

// NewICM20948Driver creates ICM20948 driver for raspberry pi
func NewICM20948Driver(settings Settings) (*Device, error) {
	d, err := sysfs.NewSPI(settings.BusNumber, settings.ChipSelect)
	if err != nil {
		return nil, err
	}
	conn, err := d.Connect(7*physic.MegaHertz, spi.Mode3, 8)
	if err != nil {
		return nil, err
	}
	dev := Device{
		Name:    "ICM20948",
		SPI:     d,
		Conn:    conn,
		regbank: 0xFF,
		acc: types.Sensor{
			Type:   ACCELEROMETER,
			Config: settings.AccConfig,
		},
		gyro: types.Sensor{
			Type:   GYROSCOPE,
			Config: settings.GyroConfig,
		},
		mag: types.Sensor{
			Type: MAGNETOMETER,
		},
	}
	return &dev, nil
}

func (dev *Device) Close() {
	dev.SPI.Close()
	fmt.Println("Closing ", dev.Name)
}

func (dev *Device) readReg(address uint8, len int) ([]uint8, error) {
	w := make([]uint8, len+1)
	r := make([]uint8, len+1)
	w[0] = (address & 0x7F) | 0x80
	err := dev.Conn.Tx(w, r)
	return r[1:], err
}

func (dev *Device) writeReg(address uint8, data ...uint8) error {
	if len(data) == 0 {
		return nil
	}
	w := append([]uint8{address & 0x7F}, data...)
	err := dev.Conn.Tx(w, nil)
	return err
}

func (dev *Device) selRegisterBank(regbank uint8) error {
	if regbank == dev.regbank {
		return nil
	}
	dev.regbank = regbank
	return dev.writeReg(REG_BANK_SEL, regbank<<4)
}

func (dev *Device) readRegister(register uint16, len int) ([]uint8, error) {
	reg := reg(register)
	dev.selRegisterBank(reg.bank)
	return dev.readReg(reg.address, len)
}

func (dev *Device) writeRegister(register uint16, data ...uint8) error {
	if len(data) == 0 {
		return nil
	}
	reg := reg(register)
	dev.selRegisterBank(reg.bank)
	return dev.writeReg(reg.address, data...)
}

// WhoAmI return value for ICM-20948 is 0xEA
func (dev *Device) WhoAmI() (name string, id uint8, err error) {
	name = "ICM-20948"
	data, err := dev.readRegister(WHO_AM_I, 1)
	id = data[0]
	return
}

// GetDeviceConfig reads device configurations
func (dev *Device) GetDeviceConfig() (
	config types.Config,
	accConfig types.Config,
	gyroConfig types.Config,
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
	// No low power mode, enabling everything with 20Mhz clock
	err = dev.writeRegister(PWR_MGMT_1, 0b00000001, 0b00000000)
	time.Sleep(50 * time.Millisecond) // wait for starting
	err = dev.InitAccelerometer()
	time.Sleep(50 * time.Millisecond) // wait for starting
	err = dev.InitGyroscope()
	time.Sleep(50 * time.Millisecond) // wait for starting
	return err
}

// ReadRawData reads all Accl and Gyro data
func (dev *Device) ReadRawData() ([]uint8, error) {
	return dev.readRegister(ACCEL_XOUT_H, 12)
}

// ResetGyroTimer resets gyro timer
func (dev *Device) ResetGyroTimer() {
	dev.lastReading = time.Now().UnixNano()
}

// ReadData reads Accelerometer and Gyro data
func (dev *Device) ReadData() (imu.ImuData, error) {
	data, err := dev.ReadRawData()
	now := time.Now().UnixNano()
	dev.duration = dev.lastReading - now
	dev.lastReading = now
	if err != nil {
		return imu.ImuData{}, err
	}
	acc, accErr := dev.processAccelerometerData(data)
	gyro, gyroErr := dev.processGyroscopeData(data[6:])

	return imu.ImuData{
		Duration: float64(dev.duration) / 1e9,
		Acc: types.SensorData{
			Error: accErr,
			Data:  acc,
		},
		Gyro: types.SensorData{
			Error: gyroErr,
			Data:  gyro,
		},
	}, nil
}

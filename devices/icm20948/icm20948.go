package icm20948

import (
	"fmt"
	"time"

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
func NewICM20948Driver(config Config) (*ImuDevice, error) {
	d, err := sysfs.NewSPI(config.BusNumber, config.ChipSelect)
	if err != nil {
		return nil, err
	}
	conn, err := d.Connect(7*physic.MegaHertz, spi.Mode3, 8)
	if err != nil {
		return nil, err
	}
	dev := ImuDevice{
		Name:    "ICM20948",
		SPI:     d,
		Conn:    conn,
		regbank: 0xFF,
		acc: types.Sensor{
			Type:   ACCELEROMETER,
			Config: config.AccConfig,
		},
		gyro: types.Sensor{
			Type:   GYROSCOPE,
			Config: config.GyroConfig,
		},
		mag: types.Sensor{
			Type: MAGNETOMETER,
		},
		prevReadTime: -1,
		readTime:     -1,
	}
	return &dev, nil
}

func (dev *ImuDevice) Close() {
	dev.SPI.Close()
	fmt.Println("Closing ", dev.Name)
}

func (dev *ImuDevice) readReg(address uint8, len int) ([]uint8, error) {
	w := make([]uint8, len+1)
	r := make([]uint8, len+1)
	w[0] = (address & 0x7F) | 0x80
	err := dev.Conn.Tx(w, r)
	return r[1:], err
}

func (dev *ImuDevice) writeReg(address uint8, data ...uint8) error {
	if len(data) == 0 {
		return nil
	}
	w := append([]uint8{address & 0x7F}, data...)
	err := dev.Conn.Tx(w, nil)
	return err
}

func (dev *ImuDevice) selRegisterBank(regbank uint8) error {
	if regbank == dev.regbank {
		return nil
	}
	dev.regbank = regbank
	return dev.writeReg(REG_BANK_SEL, regbank<<4)
}

func (dev *ImuDevice) readRegister(register uint16, len int) ([]uint8, error) {
	reg := reg(register)
	dev.selRegisterBank(reg.bank)
	return dev.readReg(reg.address, len)
}

func (dev *ImuDevice) writeRegister(register uint16, data ...uint8) error {
	if len(data) == 0 {
		return nil
	}
	reg := reg(register)
	dev.selRegisterBank(reg.bank)
	return dev.writeReg(reg.address, data...)
}

// WhoAmI return value for ICM-20948 is 0xEA
func (dev *ImuDevice) WhoAmI() (name string, id uint8, err error) {
	name = "ICM-20948"
	data, err := dev.readRegister(WHO_AM_I, 1)
	id = data[0]
	return
}

// GetDeviceConfig reads device configurations
func (dev *ImuDevice) GetDeviceConfig() (
	config types.Config,
	accConfig types.Config,
	gyroConfig types.Config,
	err error) {
	// data, err := dev.readRegister(LP_CONFIG, 3)
	accConfig, err = dev.getAccConfig()
	gyroConfig, err = dev.getGyroConfig()
	return
}

// InitDevice applies initial configurations to device
func (dev *ImuDevice) InitDevice() error {
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

// ReadSensorsRawData reads all Accl and Gyro data
func (dev *ImuDevice) ReadSensorsRawData() ([]uint8, error) {
	now := time.Now().UnixNano()
	if dev.readTime < 0 {
		dev.prevReadTime = now
	} else {
		dev.prevReadTime = dev.readTime
	}
	dev.readTime = now
	return dev.readRegister(ACCEL_XOUT_H, 12)
}

// ReadSensors reads Accelerometer and Gyro data
func (dev *ImuDevice) ReadSensors() (types.ImuSensorsData, error) {
	data, err := dev.ReadSensorsRawData()

	if err != nil {
		return types.ImuSensorsData{}, err
	}
	acc, accErr := dev.processAccelerometerData(data)
	gyro, gyroErr := dev.processGyroscopeData(data[6:])

	return types.ImuSensorsData{
		Acc: types.SensorData{
			Error: accErr,
			Data:  acc,
		},
		Gyro: types.SensorData{
			Error: gyroErr,
			Data:  gyro,
		},
		Mag: types.SensorData{
			Error: nil,
			Data:  types.XYZ{X: 0, Y: 0, Z: 0},
		},
		ReadTime:     dev.readTime,
		ReadInterval: dev.readTime - dev.prevReadTime,
	}, nil
}

func (dev *ImuDevice) GetRotations() (types.ImuRotations, error) {
	imuData, err := dev.ReadSensors()
	return types.ImuRotations{
		Accelerometer: types.Rotations{Roll: 0, Pitch: 0, Yaw: 0},
		Gyroscope:     types.Rotations{Roll: 0, Pitch: 0, Yaw: 0},
		Rotations:     types.Rotations{Roll: 0, Pitch: 0, Yaw: 0},
		ReadTime:      imuData.ReadTime,
		ReadInterval:  imuData.ReadInterval,
	}, err
}

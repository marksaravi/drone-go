package icm20948

import (
	"log"
	"time"

	"github.com/MarkSaravi/drone-go/models"
	"periph.io/x/periph/conn/spi"
)

type register struct {
	Address uint8
	Bank    uint8
}

type sensorConfig struct {
	sensitivityLevel     string
	averaging            int
	lowPassFilterEnabled bool
	lowPassFilterConfig  int
	offsetX              int16
	offsetY              int16
	offsetZ              int16
}

type memsICM20948 struct {
	spiConn    spi.Conn
	regbank    uint8
	accConfig  sensorConfig
	gyroConfig sensorConfig
	accData    models.XYZ
	gyroData   models.XYZ
}

func reg(reg uint16) *register {
	return &register{
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
func NewICM20948Driver(
	spiConn spi.Conn,
	accSensitivityLevel string,
	accAveraging int,
	accLowPassFilterEnabled bool,
	accLowPassFilterConfig int,
	accOffsetX int16,
	accOffsetY int16,
	accOffsetZ int16,
	gyroSensitivityLevel string,
	gyroAveraging int,
	gyroLowPassFilterEnabled bool,
	gyroLowPassFilterConfig int,
	gyroOffsetX int16,
	gyroOffsetY int16,
	gyroOffsetZ int16,
) *memsICM20948 {
	dev := memsICM20948{
		spiConn: spiConn,
		regbank: 0xFF,
		accConfig: sensorConfig{
			sensitivityLevel:     accSensitivityLevel,
			averaging:            accAveraging,
			lowPassFilterEnabled: accLowPassFilterEnabled,
			lowPassFilterConfig:  accLowPassFilterConfig,
			offsetX:              accOffsetX,
			offsetY:              accOffsetY,
			offsetZ:              accOffsetZ,
		},
		gyroConfig: sensorConfig{
			sensitivityLevel:     gyroSensitivityLevel,
			averaging:            gyroAveraging,
			lowPassFilterEnabled: gyroLowPassFilterEnabled,
			lowPassFilterConfig:  gyroLowPassFilterConfig,
			offsetX:              gyroOffsetX,
			offsetY:              gyroOffsetY,
			offsetZ:              gyroOffsetZ,
		},
	}
	err := dev.initDevice()
	if err != nil {
		log.Fatal(err)
	}
	return &dev
}

func (dev *memsICM20948) readReg(address uint8, len int) ([]uint8, error) {
	w := make([]uint8, len+1)
	r := make([]uint8, len+1)
	w[0] = (address & 0x7F) | 0x80
	err := dev.spiConn.Tx(w, r)
	return r[1:], err
}

func (dev *memsICM20948) writeReg(address uint8, data ...uint8) error {
	if len(data) == 0 {
		return nil
	}
	w := append([]uint8{address & 0x7F}, data...)
	err := dev.spiConn.Tx(w, nil)
	return err
}

func (dev *memsICM20948) selRegisterBank(regbank uint8) error {
	if regbank == dev.regbank {
		return nil
	}
	dev.regbank = regbank
	return dev.writeReg(REG_BANK_SEL, regbank<<4)
}

func (dev *memsICM20948) readRegister(address uint16, len int) ([]uint8, error) {
	reg := reg(address)
	dev.selRegisterBank(reg.Bank)
	return dev.readReg(reg.Address, len)
}

func (dev *memsICM20948) writeRegister(address uint16, data ...uint8) error {
	if len(data) == 0 {
		return nil
	}
	reg := reg(address)
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
	err = dev.initAccelerometer()
	if err != nil {
		return err
	}
	time.Sleep(50 * time.Millisecond) // wait for starting
	err = dev.initGyroscope()
	time.Sleep(50 * time.Millisecond) // wait for starting
	return err
}

// readSensorsRawData reads all Accl and Gyro data
func (dev *memsICM20948) readSensorsRawData() ([]uint8, error) {
	return dev.readRegister(ACCEL_XOUT_H, 12)
}

// ReadSensors reads Accelerometer and Gyro data
func (dev *memsICM20948) Read() (
	acc models.XYZ,
	gyro models.XYZ,
) {
	data, err := dev.readSensorsRawData()

	if err != nil {
		return
	}
	nacc, accErr := dev.processAccelerometerData(data)
	ngyro, gyroErr := dev.processGyroscopeData(data[6:])
	if accErr == nil {
		dev.accData = nacc
	}
	if gyroErr == nil {
		dev.gyroData = ngyro
	}
	return dev.accData, dev.gyroData
}

// towsComplementUint8ToInt16 converts 2's complement H and L uint8 to signed int16
func towsComplementUint8ToInt16(h, l uint8) int16 {
	var h16 uint16 = uint16(h)
	var l16 uint16 = uint16(l)

	return int16((h16 << 8) | l16)
}

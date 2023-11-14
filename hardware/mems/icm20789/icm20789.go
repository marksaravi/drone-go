package icm20789

import (
	"time"

	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/mems"
	"periph.io/x/conn/v3/spi"
)

const (
	ADRESS_GYRO_CONFIG   byte = 0x1B
	ADDRESS_PWR_MGMT_1   byte = 0x6B
	ADDRESS_PWR_MGMT_2   byte = 0x6C
	ADDRESS_WHO_AM_I     byte = 0x75
)

const (
	PWR_MGMT_1_CONFIG_DEVICE_RESET byte = 0b10000000
	PWR_MGMT_1_CONFIG              byte = 0b00000000
	PWR_MGMT_2_CONFIG              byte = 0b00000000
)

const (
	GYRO_FULL_SCALE_250DPS  float64 = 131
	GYRO_FULL_SCALE_500DPS  float64 = 65.5
	GYRO_FULL_SCALE_1000DPS float64 = 32.8
	GYRO_FULL_SCALE_2000DPS float64 = 16.4
)

const (
	RAW_DATA_SIZE int = 14
)

type Offsets struct {
	X uint16 `json:"x"`
	Y uint16 `json:"y"`
	Z uint16 `json:"z"`
}

type InertialDeviceConfigs struct {
	FullScale string  `json:"full_scale"`
	Offsets   Offsets `json:"offsets"`
}

type Configs struct {
	SPI           hardware.SPIConnConfigs `json:"spi"`
	Accelerometer InertialDeviceConfigs   `json:"accelerometer"`
	Gyroscope     InertialDeviceConfigs   `json:"gyroscope"`
}

type memsIcm20789 struct {
	spiConn spi.Conn

	accelFullScale float64
	gyroFullScale  float64
}

func NewICM20789(configs Configs) *memsIcm20789 {
	gyroFullScale, gyroFullScaleMask := gyroscopeFullScale(configs.Gyroscope.FullScale)
	m := memsIcm20789{
		spiConn:        hardware.NewSPIConnection(configs.SPI),
		accelFullScale: ACCEL_FULL_SCALE_G[configs.Accelerometer.FullScale],
		gyroFullScale:  gyroFullScale,
	}
	m.setupPower()
	// (fullScale string, numberOfSamples int, fifoSize int, lowPassFilterFrequency string)
	m.setupAccelerometer(configs.Accelerometer.FullScale, 8, 512, "44.8hz")
	m.setupGyroscope(gyroFullScaleMask)
	return &m
}

func (m *memsIcm20789) WhoAmI() (byte, error) {
	memsData, err := m.readRegister(ADDRESS_WHO_AM_I, 1)
	return memsData[0], err
}

func (m *memsIcm20789) Read() (mems.Mems6DOFData, error) {
	memsData, err := m.readRegister(ADDRESS_ACCEL_XOUT_H, RAW_DATA_SIZE)
	if err != nil {
		return mems.Mems6DOFData{}, err
	}
	return mems.Mems6DOFData{
		Accelerometer: m.memsDataToAccelerometer(memsData[:6]),
		Gyroscope:     m.memsDataToGyroscope(memsData[8:]), // 6 and 7 are Temperature data
	}, nil
}

func (m *memsIcm20789) readRegister(address byte, size int) ([]byte, error) {
	w := make([]byte, size+1)
	r := make([]byte, size+1)
	w[0] = address | byte(0x80)

	err := m.spiConn.Tx(w, r)
	return r[1:], err
}

func (m *memsIcm20789) readByteFromRegister(address byte) (byte, error) {
	res, err := m.readRegister(address, 1)
	return res[0], err
}

func (m *memsIcm20789) writeRegister(address byte, data ...byte) error {
	w := make([]byte, 1, len(data)+1)
	r := make([]byte, cap(w))
	w[0] = address
	w = append(w, data...)
	err := m.spiConn.Tx(w, r)
	return err
}

func (m *memsIcm20789) setupPower() {
	m.writeRegister(ADDRESS_PWR_MGMT_1, PWR_MGMT_1_CONFIG_DEVICE_RESET)
	delay(10)
	m.writeRegister(ADDRESS_PWR_MGMT_1, PWR_MGMT_1_CONFIG)
	delay(10)
}

func (m *memsIcm20789) memsDataToGyroscope(memsData []byte) mems.DXYZ {
	return mems.DXYZ{
		DX: float64(towsComplementUint8ToInt16(memsData[0], memsData[1])) / m.gyroFullScale,
		DY: float64(towsComplementUint8ToInt16(memsData[2], memsData[3])) / m.gyroFullScale,
		DZ: float64(towsComplementUint8ToInt16(memsData[4], memsData[5])) / m.gyroFullScale,
	}
}

// towsComplementUint8ToInt16 converts 2's complement H and L uint8 to signed int16
func towsComplementUint8ToInt16(h, l byte) int16 {
	h16 := uint16(h)
	l16 := uint16(l)

	return int16(h16<<8 | l16)
}

func delay(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

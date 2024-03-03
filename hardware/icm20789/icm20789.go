package icm20789

import (
	"log"
	"time"

	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/mems"
	"periph.io/x/conn/v3/spi"
)

const (
	ADDRESS_PWR_MGMT_1 byte = 0x6B
	ADDRESS_PWR_MGMT_2 byte = 0x6C
	ADDRESS_WHO_AM_I   byte = 0x75
)

const (
	PWR_MGMT_1_CONFIG_DEVICE_RESET byte = 0b10000000
	PWR_MGMT_1_CONFIG              byte = 0b00000000
	PWR_MGMT_2_CONFIG              byte = 0b00000000
)

const (
	RAW_DATA_SIZE int = 14
)

type Offsets struct {
	X uint16 `json:"x"`
	Y uint16 `json:"y"`
	Z uint16 `json:"z"`
}

type AccelerometerConfigs struct {
	FullScale              string  `json:"full_scale"`
	Offsets                Offsets `json:"offsets"`
	LowPassFilterFrequency string  `json:"lowpass_filter_frequency"`
	NumberOfSamples        int     `json:"number_of_samples"`
}

type GyroscopeConfigs struct {
	FullScale string  `json:"full-scale"`
	Offsets   Offsets `json:"offsets"`
}

type Configs struct {
	Accelerometer AccelerometerConfigs    `json:"accelerometer"`
	Gyroscope     GyroscopeConfigs        `json:"gyroscope"`
	SPI           hardware.SPIConnConfigs `json:"spi"`
}

type memsIcm20789 struct {
	spiConn spi.Conn

	accelFullScale float64
	gyroFullScale  float64
}

func NewICM20789(configs Configs) *memsIcm20789 {
	var accelFullScale, gyroFullScale float64
	var ok bool
	if accelFullScale, ok = ACCEL_FULL_SCALE_G[configs.Accelerometer.FullScale]; !ok {
		log.Fatalf("Error: Accelerometer Full Scale is not defined.")
	}
	if gyroFullScale, ok = GYRO_FULL_SCALE_DPS[configs.Gyroscope.FullScale]; !ok {
		log.Fatalf("Error: Gyroscope Full Scale is not defined.")
	}
	log.Println("GYROSCOPE FULL SCALE: ", gyroFullScale)
	m := memsIcm20789{
		spiConn:        hardware.NewSPIConnection(configs.SPI),
		accelFullScale: accelFullScale,
		gyroFullScale:  gyroFullScale,
	}
	m.setupPower()
	m.setupAccelerometer(
		configs.Accelerometer.FullScale,
		configs.Accelerometer.NumberOfSamples,
		512,
		configs.Accelerometer.LowPassFilterFrequency,
		configs.Accelerometer.Offsets.X,
		configs.Accelerometer.Offsets.Y,
		configs.Accelerometer.Offsets.Z,
	)
	m.setupGyroscope(
		configs.Gyroscope.FullScale,
		configs.Gyroscope.Offsets.X,
		configs.Gyroscope.Offsets.Y,
		configs.Gyroscope.Offsets.Z,
	)
	return &m
}

func (m *memsIcm20789) WhoAmI() (byte, error) {
	memsData, err := m.readRegister(ADDRESS_WHO_AM_I, 1)
	return memsData[0], err
}

func (m *memsIcm20789) Read() (mems.Mems6DOFData, error) {
	memsData, err := m.readRegister(DATA_READ_SEGMENT, RAW_DATA_SIZE)
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

// towsComplementUint8ToInt16 converts 2's complement H and L uint8 to signed int16
func towsComplementUint8ToInt16(h, l byte) int16 {
	h16 := uint16(h)
	l16 := uint16(l)

	return int16(h16<<8 | l16)
}

func delay(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

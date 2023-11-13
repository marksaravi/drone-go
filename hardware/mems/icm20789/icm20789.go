package icm20789

import (
	"time"

	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/mems"
	"periph.io/x/conn/v3/spi"
)

const (
	ADDRESS_ACCEL_XOUT_H byte = 0x3B
	ADDRESS_ACCEL_CONFIG byte = 0x1C
	ADRESS_GYRO_CONFIG   byte = 0x1B
	ADDRESS_PWR_MGMT_1   byte = 0x6B
	ADDRESS_PWR_MGMT_2   byte = 0x6C
	ADDRESS_WHO_AM_I     byte = 0x75
	ADDRESS_XA_OFFSH     byte = 0x77
	ADDRESS_XA_OFFSL     byte = 0x78
	ADDRESS_YA_OFFSH     byte = 0x7A
	ADDRESS_YA_OFFSL     byte = 0x7B
	ADDRESS_ZA_OFFSH     byte = 0x7D
	ADDRESS_ZA_OFFSL     byte = 0x7E
)

const (
	PWR_MGMT_1_CONFIG_DEVICE_RESET byte = 0b10000000
	PWR_MGMT_1_CONFIG              byte = 0b00000000
	PWR_MGMT_2_CONFIG              byte = 0b00000000
)

const (
	ACCEL_CONFIG_DISABLE_SELF_TESTS  byte = 0b00000000
	ACCEL_CONFIG_2G                  byte = 0b00000000
	ACCEL_CONFIG_4G                  byte = 0b00001000
	ACCEL_CONFIG_8G                  byte = 0b00010000
	ACCEL_CONFIG_16G                 byte = 0b00011000
	ACCEL_CONFIG2_FIFO_SIZE_512      byte = 0b00000000
	ACCEL_CONFIG2_DEC2_CFG_4_SAMPLE  byte = 0b00000000
	ACCEL_CONFIG2_DEC2_CFG_8_SAMPLE  byte = 0b00010000
	ACCEL_CONFIG2_DEC2_CFG_16_SAMPLE byte = 0b00100000
	ACCEL_CONFIG2_DEC2_CFG_32_SAMPLE byte = 0b00110000
	ACCEL_CONFIG2_ACCEL_FCHOICE_B    byte = 0b00001000 //3-dB BW (Hz) 1046.0
	ACCEL_CONFIG2_A_DLPF_CFG_3dB_BW0 byte = 0b00000000 //3-dB BW (Hz) 218.1
	ACCEL_CONFIG2_A_DLPF_CFG_3dB_BW1 byte = 0b00000001 //3-dB BW (Hz) 218.1
	ACCEL_CONFIG2_A_DLPF_CFG_3dB_BW2 byte = 0b00000010 //3-dB BW (Hz) 99.0
	ACCEL_CONFIG2_A_DLPF_CFG_3dB_BW3 byte = 0b00000011 //3-dB BW (Hz) 44.8
	ACCEL_CONFIG2_A_DLPF_CFG_3dB_BW4 byte = 0b00000100 //3-dB BW (Hz) 21.2
	ACCEL_CONFIG2_A_DLPF_CFG_3dB_BW5 byte = 0b00000101 //3-dB BW (Hz) 10.2
	ACCEL_CONFIG2_A_DLPF_CFG_3dB_BW6 byte = 0b00000110 //3-dB BW (Hz) 5.1
	ACCEL_CONFIG2_A_DLPF_CFG_3dB_BW7 byte = 0b00000111 //3-dB BW (Hz) 420.0
)

const (
	ACCEL_FULL_SCALE_2G  float64 = 16384
	ACCEL_FULL_SCALE_4G  float64 = 8192
	ACCEL_FULL_SCALE_8G  float64 = 4096
	ACCEL_FULL_SCALE_16G float64 = 2048

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
	accelFullScale, accelFullScaleMask := accelerometerFullScale(configs.Accelerometer.FullScale)
	gyroFullScale, gyroFullScaleMask := gyroscopeFullScale(configs.Gyroscope.FullScale)
	m := memsIcm20789{
		spiConn:        hardware.NewSPIConnection(configs.SPI),
		accelFullScale: accelFullScale,
		gyroFullScale:  gyroFullScale,
	}
	m.setupPower()
	m.setupAccelerometer(accelFullScaleMask)
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
		Accelerometer: m.memsDataToAccelerometer(memsData),
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

func (m *memsIcm20789) memsDataToAccelerometer(memsData []byte) mems.XYZ {
	return mems.XYZ{
		X: float64(towsComplementUint8ToInt16(memsData[0], memsData[1])) / m.accelFullScale,
		Y: float64(towsComplementUint8ToInt16(memsData[2], memsData[3])) / m.accelFullScale,
		Z: float64(towsComplementUint8ToInt16(memsData[4], memsData[5])) / m.accelFullScale,
	}
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

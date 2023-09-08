package icm20789

import (
	"time"

	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/mems"
	"github.com/marksaravi/drone-go/utils"
	"periph.io/x/conn/v3/spi"
)

const (
	ACCEL_XOUT_H byte = 0x3B
	ACCEL_CONFIG byte = 0x1C
	GYRO_CONFIG  byte = 0x1B
	WHO_AM_I     byte = 0x75
	PWR_MGMT_1   byte = 0x6B
	PWR_MGMT_2   byte = 0x6C
)

const (
	PWR_MGMT_1_CONFIG byte = 0b00000000
	PWR_MGMT_2_CONFIG byte = 0b00000000
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
	X uint16 `yaml:"x"`
	Y uint16 `yaml:"y"`
	Z uint16 `yaml:"z"`
}

type InertialDeviceConfigs struct {
	FullScale string  `yaml:"full_scale"`
	Offsets   Offsets `yaml:"offsets"`
}

type Configs struct {
	Accelerometer InertialDeviceConfigs `yaml:"accelerometer"`
	Gyroscope     InertialDeviceConfigs `yaml:"gyroscope"`
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
		spiConn:        hardware.NewSPIConnection(0, 0),
		accelFullScale: accelFullScale,
		gyroFullScale:  gyroFullScale,
	}
	m.setupPower()
	m.setupAccelerometer(accelFullScaleMask)
	m.setupGyroscope(gyroFullScaleMask)
	return &m
}

func ReadConfigs() Configs {
	var configs struct {
		Imu Configs `yaml:"icm20789"`
	}
	utils.ReadConfigs(&configs)
	return configs.Imu
}
func (m *memsIcm20789) Read() (mems.Mems6DOFData, error) {
	memsData, err := m.readRegister(ACCEL_XOUT_H, RAW_DATA_SIZE)
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
	m.writeRegister(PWR_MGMT_1, 0x80) // soft reset
	delay(1)
	m.writeRegister(PWR_MGMT_1, PWR_MGMT_1_CONFIG)
	delay(1)
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
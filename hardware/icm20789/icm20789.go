package icm20789

import (
	"time"

	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/types"
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

type ICM20789Configs struct {
	AccelerometerFullScale string `yaml:"accelerometer_full_scale"`
	AccelerometerXOffset   uint16 `yaml:"accelerometer_x_offset"`
	AccelerometerYOffset   uint16 `yaml:"accelerometer_y_offset"`
	AccelerometerZOffset   uint16 `yaml:"accelerometer_z_offset"`

	GyroscopeFullScale string `yaml:"gyroscope_full_scale"`
	GyroscopeXOffset   uint16 `yaml:"gyroscope_x_offset"`
	GyroscopeYOffset   uint16 `yaml:"gyroscope_y_offset"`
	GyroscopeZOffset   uint16 `yaml:"gyroscope_z_offset"`
}

type imuICM20789 struct {
	spiConn spi.Conn

	accelFullScale float64
	gyroFullScale  float64
}

func NewICM20789(configs ICM20789Configs) *imuICM20789 {
	accelFullScale, accelFullScaleMask := accelerometerFullScale(configs.AccelerometerFullScale)
	gyroFullScale, gyroFullScaleMask := gyroscopeFullScale(configs.GyroscopeFullScale)
	imu := imuICM20789{
		spiConn:        hardware.NewSPIConnection(0, 0),
		accelFullScale: accelFullScale,
		gyroFullScale:  gyroFullScale,
	}
	imu.setupPower()
	imu.setupAccelerometer(accelFullScaleMask)
	imu.setupGyroscope(gyroFullScaleMask)

	return &imu
}

func (imu *imuICM20789) Read() (types.IMUMems6DOFRawData, error) {
	memsData, err := imu.readRegister(ACCEL_XOUT_H, RAW_DATA_SIZE)
	if err != nil {
		return types.IMUMems6DOFRawData{}, err
	}
	return types.IMUMems6DOFRawData{
		Accelerometer: imu.memsDataToAccelerometer(memsData),
		Gyroscope:     imu.memsDataToGyroscope(memsData[8:]), // 6 and 7 are Temperature data
	}, nil
}

func (imu *imuICM20789) readRegister(address byte, size int) ([]byte, error) {
	w := make([]byte, size+1)
	r := make([]byte, size+1)
	w[0] = address | byte(0x80)

	err := imu.spiConn.Tx(w, r)
	return r[1:], err
}

func (imu *imuICM20789) readByteFromRegister(address byte) (byte, error) {
	res, err := imu.readRegister(address, 1)
	return res[0], err
}

func (imu *imuICM20789) writeRegister(address byte, data ...byte) error {
	w := make([]byte, 1, len(data)+1)
	r := make([]byte, cap(w))
	w[0] = address
	w = append(w, data...)
	err := imu.spiConn.Tx(w, r)
	return err
}

func (imu *imuICM20789) setupPower() {
	imu.writeRegister(PWR_MGMT_1, 0x80) // soft reset
	delay(1)
	imu.writeRegister(PWR_MGMT_1, PWR_MGMT_1_CONFIG)
	delay(1)
}

func (imu *imuICM20789) memsDataToAccelerometer(memsData []byte) types.XYZ {
	return types.XYZ{
		X: float64(towsComplementUint8ToInt16(memsData[0], memsData[1])) / imu.accelFullScale,
		Y: float64(towsComplementUint8ToInt16(memsData[2], memsData[3])) / imu.accelFullScale,
		Z: float64(towsComplementUint8ToInt16(memsData[4], memsData[5])) / imu.accelFullScale,
	}
}

func (imu *imuICM20789) memsDataToGyroscope(memsData []byte) types.DXYZ {
	return types.DXYZ{
		DX: float64(towsComplementUint8ToInt16(memsData[0], memsData[1])) / imu.gyroFullScale,
		DY: float64(towsComplementUint8ToInt16(memsData[2], memsData[3])) / imu.gyroFullScale,
		DZ: float64(towsComplementUint8ToInt16(memsData[4], memsData[5])) / imu.gyroFullScale,
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

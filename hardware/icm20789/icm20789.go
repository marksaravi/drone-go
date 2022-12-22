package icm20789

import (
	"fmt"
	"time"

	"github.com/marksaravi/drone-go/hardware"
)

type spi interface{}

const (
	ACCELEROMETER_DATA_SIZE = 6
	GYROSCOPE_DATA_SIZE     = ACCELEROMETER_DATA_SIZE
)

const (
	ACCEL_CONFIG byte = 0x1C
	ACCEL_XOUT_H byte = 0x3B
	GYRO_CONFIG  byte = 0x1B
	WHO_AM_I     byte = 0x75
	PWR_MGMT_1   byte = 0x6B
	XA_OFFSET_H  byte = 0x13
)

const (
	ACCEL_CONFIG_MASK_FULL_SCALE_2G     byte = 0b00000000
	ACCEL_CONFIG_MASK_FULL_SCALE_4G     byte = 0b00001000
	ACCEL_CONFIG_MASK_FULL_SCALE_8G     byte = 0b00010000
	ACCEL_CONFIG_MASK_FULL_SCALE_16G    byte = 0b00011000
	GYRO_CONFIG_MASK_FULL_SCALE_250DPS  byte = 0b00000000
	GYRO_CONFIG_MASK_FULL_SCALE_500DPS  byte = 0b00001000
	GYRO_CONFIG_MASK_FULL_SCALE_1000DPS byte = 0b00010000
	GYRO_CONFIG_MASK_FULL_SCALE_2000DPS byte = 0b00011000
)

type imuIcm20789 struct {
	spiConn spi.Conn
}

func NewIcm20987(spiBusNumber, spiChiSelect int) *imuIcm20789 {
	spiConn := hardware.NewSPIConnection(spiBusNumber, spiChiSelect)
	return &imuIcm20789{
		spiConn: spiConn,
	}
}

func (imu *imuIcm20789) readRegister(address byte, size int) ([]byte, error) {
	w := make([]byte, size+1)
	r := make([]byte, size+1)
	w[0] = address | byte(0x80)

	err := imu.spiConn.Tx(w, r)
	return r[1:], err
}

func (imu *imuIcm20789) readByteFromRegister(address byte) (byte, error) {
	res, err := imu.readRegister(address, 1)
	return res[0], err
}

func (imu *imuIcm20789) writeRegister(address byte, data ...byte) error {
	w := make([]byte, 1, len(data)+1)
	r := make([]byte, cap(w))
	w[0] = address
	w = append(w, data...)
	fmt.Println("len: ", len(w), len(r))
	err := imu.spiConn.Tx(w, r)
	return err
}

func (imu *imuIcm20789) SPIReadTest() (whoami, powermanagement1 byte, ok bool, err error) {
	whoami, err = imu.readByteFromRegister(WHO_AM_I)
	if err == nil {
		powermanagement1, err = imu.readByteFromRegister(PWR_MGMT_1)
	}
	return whoami, powermanagement1, whoami == 0x03 && powermanagement1 == 0x40 && err == nil, err
}
func (imu *imuIcm20789) SPIWriteTest() (bool, byte, byte, error) {
	v, err := imu.readRegister(XA_OFFSET_H, 1)
	if err != nil {
		return false, v[0], 0, err
	}
	err = imu.writeRegister(XA_OFFSET_H, v[0]+5)
	if err != nil {
		return false, v[0], 0, err
	}
	nv, err := imu.readRegister(XA_OFFSET_H, 1)
	return nv[0] == v[0]+5, v[0], nv[0], err
}

func (imu *imuIcm20789) ReadAccelerometer() ([]byte, error) {
	return imu.readRegister(107, 1)
}

func delay(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

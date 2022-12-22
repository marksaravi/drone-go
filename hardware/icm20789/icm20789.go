package icm20789

import (
	"time"

	"github.com/marksaravi/drone-go/hardware"
	"periph.io/x/periph/conn/spi"
)

const (
	ACCELEROMETER_DATA_SIZE = 6
	GYROSCOPE_DATA_SIZE     = ACCELEROMETER_DATA_SIZE
)

const (
	ACCEL_CONFIG byte = 0x1C
	ACCEL_XOUT_H byte = 0x3B
	GYRO_CONFIG  byte = 0x1B
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

func (imu *imuIcm20789) readSPI(address byte, size int) ([]byte, error) {
	w := make([]byte, size+1)
	r := make([]byte, size+1)
	w[0] = address | byte(0x80)

	err := imu.spiConn.Tx(w, r)
	return r[1:], err
}

func (imu *imuIcm20789) writeSPI(address byte, data []byte) error {
	w := make([]byte, 0, len(data)+1)

	w = append(w, address)
	w = append(w, data...)
	r := make([]byte, len(w))

	err := imu.spiConn.Tx(w, r)
	return err
}

func (imu *imuIcm20789) WhoAmI() (byte, error) {
	data, err := imu.readSPI(0x75, 1)
	return data[0], err
}

func (imu *imuIcm20789) ReadAccelerometer() ([]byte, error) {
	return imu.readSPI(ACCEL_XOUT_H, ACCELEROMETER_DATA_SIZE)
}

func delay(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

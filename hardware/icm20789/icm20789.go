package icm20789

import (
	"fmt"

	"github.com/marksaravi/drone-go/hardware"
	"periph.io/x/periph/conn/spi"
)

const (
	ACCELEROMETER_DATA_SIZE = 6
	GYROSCOPE_DATA_SIZE     = ACCELEROMETER_DATA_SIZE
)

const (
	ACCEL_XOUT_H uint8 = 0x3B
)

type Settings struct {
	SPI hardware.SPISettings
}

type imuIcm20789 struct {
	spiConn spi.Conn
}

func NewIcm20987(spiBusNumber, spiChiSelect int) *imuIcm20789 {
	spiConn := hardware.NewSPIConnection(spiBusNumber, spiChiSelect)
	return &imuIcm20789{
		spiConn: spiConn,
	}
}

func (imu *imuIcm20789) readSPI(address uint8, size int) ([]uint8, error) {
	w := make([]uint8, size+1)
	r := make([]uint8, size+1)
	w[0] = address | uint8(0x80)

	err := imu.spiConn.Tx(w, r)
	fmt.Println(w, r)
	return r[1:], err
}

func (imu *imuIcm20789) ReadAccelerometer() ([]uint8, error) {
	return imu.readSPI(ACCEL_XOUT_H, ACCELEROMETER_DATA_SIZE)
}

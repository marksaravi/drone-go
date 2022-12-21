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
	imu := imuIcm20789{
		spiConn: spiConn,
	}
	imu.initialise()
	return &imu
}

func (imu *imuIcm20789) readSPI(address uint8, size int) ([]uint8, error) {
	w := make([]uint8, size+1)
	r := make([]uint8, size+1)
	w[0] = address | uint8(0x80)

	err := imu.spiConn.Tx(w, r)
	return r[1:], err
}

func (imu *imuIcm20789) writeSPI(address uint8, data []uint8) error {
	w := make([]uint8, 0, len(data)+1)

	w = append(w, address)
	w = append(w, data...)

	err := imu.spiConn.Tx(w, nil)
	return err
}

func (imu *imuIcm20789) initialise() {
	const DELAY = 10
	imu.writeSPI(0x6B, []uint8{0x01})
	imu.writeSPI(0x6B, []uint8{0x01})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x6A, []uint8{0x10})
	imu.writeSPI(0x6B, []uint8{0x41})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x6C, []uint8{0x3f})
	imu.writeSPI(0xF5, []uint8{0x00})
	imu.writeSPI(0x19, []uint8{0x09})
	imu.writeSPI(0xEA, []uint8{0x00})
	imu.writeSPI(0x6B, []uint8{0x01})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x6A, []uint8{0x10})
	imu.writeSPI(0x6B, []uint8{0x41})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x6B, []uint8{0x01})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x23, []uint8{0x00})
	imu.writeSPI(0x6B, []uint8{0x41})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x1D, []uint8{0xC0})
	imu.writeSPI(0x6B, []uint8{0x01})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x1A, []uint8{0xC0})
	imu.writeSPI(0x6B, []uint8{0x41})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x38, []uint8{0x01})

	// spi_dev->read_registers(0x6A, &v, 1);
	// printf("reg 0x6A=0x%02x\n", v);
	// spi_dev->read_registers(0x6B, &v, 1);
	// printf("reg 0x6B=0x%02x\n", v);

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x6A, []uint8{0x10})
	imu.writeSPI(0x6B, []uint8{0x41})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x6B, []uint8{0x01})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x23, []uint8{0x00})
	imu.writeSPI(0x6B, []uint8{0x41})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x6B, []uint8{0x41})
	imu.writeSPI(0x6C, []uint8{0x3f})
	imu.writeSPI(0x6B, []uint8{0x41})

	// spi_dev->read_registers(0x6A, &v, 1);
	// printf("reg 0x6A=0x%02x\n", v);
	imu.writeSPI(0x6B, []uint8{0x01})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x6A, []uint8{0x10})
	imu.writeSPI(0x6B, []uint8{0x41})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x6B, []uint8{0x01})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x23, []uint8{0x00})
	imu.writeSPI(0x6B, []uint8{0x41})

	time.Sleep(DELAY * time.Millisecond)
	// spi_dev->read_registers(0x6A, &v, 1);
	// printf("reg 0x6A=0x%02x\n", v);
	imu.writeSPI(0x6B, []uint8{0x01})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x6A, []uint8{0x10})
	imu.writeSPI(0x6B, []uint8{0x41})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x6B, []uint8{0x01})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x23, []uint8{0x00})
	imu.writeSPI(0x6B, []uint8{0x41})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x6B, []uint8{0x41})
	imu.writeSPI(0x6C, []uint8{0x3f})
	imu.writeSPI(0x6B, []uint8{0x41})

	// spi_dev->read_registers(0x6A, &v, 1);
	// printf("reg 0x6A=0x%02x\n", v);
	imu.writeSPI(0x6B, []uint8{0x01})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x6A, []uint8{0x10})
	imu.writeSPI(0x6B, []uint8{0x41})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x6B, []uint8{0x01})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x23, []uint8{0x00})
	imu.writeSPI(0x6B, []uint8{0x41})

}

func (imu *imuIcm20789) ReadAccelerometer() ([]uint8, error) {
	return imu.readSPI(ACCEL_XOUT_H, ACCELEROMETER_DATA_SIZE)
}

func (imu *imuIcm20789) WhoAmI() ([]uint8, error) {
	return imu.readSPI(0x75, 1)
}

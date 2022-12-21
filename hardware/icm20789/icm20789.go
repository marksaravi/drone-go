package icm20789

import (
	"log"
	"time"

	"github.com/marksaravi/drone-go/hardware"
	"periph.io/x/periph/conn/spi"
)

const (
	ACCELEROMETER_DATA_SIZE = 6
	GYROSCOPE_DATA_SIZE     = ACCELEROMETER_DATA_SIZE
)

const (
	ACCEL_XOUT_H     byte = 0x3B
	GYROSCOPE_CONFIG byte = 0x1B
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

	err := imu.spiConn.Tx(w, nil)
	return err
}

func (imu *imuIcm20789) WhoAmI() ([]byte, error) {
	return imu.readSPI(0x75, 1)
}

func (imu *imuIcm20789) Initialize() {
	const DELAY = 10
	imu.writeSPI(0x6B, []byte{0x01})
	imu.writeSPI(0x6B, []byte{0x01})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x6A, []byte{0x10})
	imu.writeSPI(0x6B, []byte{0x41})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x6C, []byte{0x3f})
	imu.writeSPI(0xF5, []byte{0x00})
	imu.writeSPI(0x19, []byte{0x09})
	imu.writeSPI(0xEA, []byte{0x00})
	imu.writeSPI(0x6B, []byte{0x01})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x6A, []byte{0x10})
	imu.writeSPI(0x6B, []byte{0x41})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x6B, []byte{0x01})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x23, []byte{0x00})
	imu.writeSPI(0x6B, []byte{0x41})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x1D, []byte{0xC0})
	imu.writeSPI(0x6B, []byte{0x01})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x1A, []byte{0xC0})
	imu.writeSPI(0x6B, []byte{0x41})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x38, []byte{0x01})

	// spi_dev->read_registers(0x6A, &v, 1);
	// printf("reg 0x6A=0x%02x\n", v);
	// spi_dev->read_registers(0x6B, &v, 1);
	// printf("reg 0x6B=0x%02x\n", v);

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x6A, []byte{0x10})
	imu.writeSPI(0x6B, []byte{0x41})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x6B, []byte{0x01})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x23, []byte{0x00})
	imu.writeSPI(0x6B, []byte{0x41})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x6B, []byte{0x41})
	imu.writeSPI(0x6C, []byte{0x3f})
	imu.writeSPI(0x6B, []byte{0x41})

	// spi_dev->read_registers(0x6A, &v, 1);
	// printf("reg 0x6A=0x%02x\n", v);
	imu.writeSPI(0x6B, []byte{0x01})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x6A, []byte{0x10})
	imu.writeSPI(0x6B, []byte{0x41})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x6B, []byte{0x01})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x23, []byte{0x00})
	imu.writeSPI(0x6B, []byte{0x41})

	time.Sleep(DELAY * time.Millisecond)
	// spi_dev->read_registers(0x6A, &v, 1);
	// printf("reg 0x6A=0x%02x\n", v);
	imu.writeSPI(0x6B, []byte{0x01})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x6A, []byte{0x10})
	imu.writeSPI(0x6B, []byte{0x41})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x6B, []byte{0x01})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x23, []byte{0x00})
	imu.writeSPI(0x6B, []byte{0x41})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x6B, []byte{0x41})
	imu.writeSPI(0x6C, []byte{0x3f})
	imu.writeSPI(0x6B, []byte{0x41})

	// spi_dev->read_registers(0x6A, &v, 1);
	// printf("reg 0x6A=0x%02x\n", v);
	imu.writeSPI(0x6B, []byte{0x01})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x6A, []byte{0x10})
	imu.writeSPI(0x6B, []byte{0x41})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x6B, []byte{0x01})

	time.Sleep(DELAY * time.Millisecond)
	imu.writeSPI(0x23, []byte{0x00})
	imu.writeSPI(0x6B, []byte{0x41})

	imu.setGyroConfigs()
}

func (imu *imuIcm20789) setGyroConfigs() {
	config, err := imu.readSPI(GYROSCOPE_CONFIG, 1)
	if err != nil {
		log.Fatalf("can't configure the gyroscope %v", err)
	}
	log.Printf("gyro: %b\n", config[0])
}

func (imu *imuIcm20789) ReadAccelerometer() ([]byte, error) {
	return imu.readSPI(ACCEL_XOUT_H, ACCELEROMETER_DATA_SIZE)
}

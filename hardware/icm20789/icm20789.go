package icm20789

import (
	"log"
	"time"

	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/types"
	"periph.io/x/conn/v3/spi"
)

const (
	ACCEL_XOUT_H byte = 0x3B
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
	RAW_DATA_SIZE int = 12
)

type imuIcm20789 struct {
	spiConn spi.Conn
}

func NewICM20789() *imuIcm20789 {
	return &imuIcm20789{
		spiConn: hardware.NewSPIConnection(0, 0),
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
	err := imu.spiConn.Tx(w, r)
	return err
}

func (imu *imuIcm20789) Setup() {
	log.Println("SETUP IMU")
	imu.setupPower()
	imu.setupGyro()
}

func (imu *imuIcm20789) setupPower() {
	log.Println("SETUP IMU power")
	imu.writeRegister(PWR_MGMT_1, 0x80) // soft reset
	delay(1)
	powerManagement1v1, _ := imu.readByteFromRegister(PWR_MGMT_1)
	imu.writeRegister(PWR_MGMT_1, PWR_MGMT_1_CONFIG)
	delay(1)
	powerManagement1v2, _ := imu.readByteFromRegister(PWR_MGMT_1)
	log.Printf("PWR_MGMT_1_v1: 0x%x, PWR_MGMT_1_v2: 0x%x\n", powerManagement1v1, powerManagement1v2)
}

func (imu *imuIcm20789) ReadIMUData() (types.IMUMems6DOFRawData, error) {
	data, err := imu.readRegister(ACCEL_XOUT_H, RAW_DATA_SIZE)
	if err != nil {
		return types.IMUMems6DOFRawData{}, err
	}
	return types.IMUMems6DOFRawData{
		Accelerometer: types.XYZ{
			X: float64(dataToFloat64(data[0], data[1])),
			Y: float64(dataToFloat64(data[2], data[3])),
			Z: float64(dataToFloat64(data[4], data[5])),
		},
		Gyroscope: types.XYZDt{
			DX: float64(dataToFloat64(data[6], data[7])),
			DY: float64(dataToFloat64(data[8], data[9])),
			DZ: float64(dataToFloat64(data[10], data[11])),
		},
	}, nil
}

func dataToFloat64(h, l byte) uint16 {
	return uint16(h)<<8 | uint16(l)
}

func delay(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

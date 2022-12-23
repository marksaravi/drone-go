package main

import (
	"fmt"
	"log"
	"time"

	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/sysfs"
)

const (
	GYRO_CONFIG byte = 0x1B
	WHO_AM_I    byte = 0x75
	PWR_MGMT_1  byte = 0x6B
	PWR_MGMT_2  byte = 0x6C
)

const (
	PWR_MGMT_1_CONFIG byte = 0b00000000
	PWR_MGMT_2_CONFIG byte = 0b00000000
)

type imuIcm20789 struct {
	spiConn spi.Conn
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

func (imu *imuIcm20789) setup() {
	imu.writeRegister(PWR_MGMT_1, 0x80) // soft reset
	delay(1)
	powerManagement1v1, _ := imu.readByteFromRegister(PWR_MGMT_1)
	imu.writeRegister(PWR_MGMT_1, PWR_MGMT_1_CONFIG)
	delay(1)
	powerManagement1v2, _ := imu.readByteFromRegister(PWR_MGMT_1)
	fmt.Printf("PWR_MGMT_1_v1: 0x%x, PWR_MGMT_1_v2: 0x%x\n", powerManagement1v1, powerManagement1v2)
	// imu.writeRegister(PWR_MGMT_1, powerManagement1)
}

func (imu *imuIcm20789) setupGyro() {
	gyrosetup1, _ := imu.readByteFromRegister(GYRO_CONFIG)
	imu.writeRegister(GYRO_CONFIG, gyrosetup1|0b00110000)
	delay(1)
	gyrosetup2, _ := imu.readByteFromRegister(GYRO_CONFIG)
	fmt.Printf("GYRO_SETUP1: 0x%x, GYRO_SETUP2: 0x%x\n", gyrosetup1, gyrosetup2)
}

func main() {
	HostInitialize()

	fmt.Println("initializing SPI")
	c := NewSPIConnection(0, 0)
	imu := imuIcm20789{
		spiConn: c,
	}
	imu.setup()
	imu.setupGyro()
	// whoami, _ := imu.readByteFromRegister(0x75)
	// fmt.Printf("WHO AM I: 0x%x\n", whoami)

	// power, _ := imu.readByteFromRegister(107)
	// fmt.Printf("POWER: 0x%x\n", power)

	// imu.writeRegister(19, 13)
}

func HostInitialize() {
	state, err := host.Init()
	if err != nil {
		log.Fatalf("failed to initialize periph: %v", err)
	}
	spiloaded := false
	i2cloaded := false

	for _, driver := range state.Loaded {
		if driver.String() == "sysfs-spi" {
			spiloaded = true
		}
		if driver.String() == "sysfs-i2c" {
			i2cloaded = true
		}
	}
	if !spiloaded {
		log.Fatalf("failed to initialize spi")
	}
	if !i2cloaded {
		log.Fatalf("failed to initialize i2c")
	}
}

func NewSPIConnection(busNumber int, chipSelect int) spi.Conn {
	p, err := sysfs.NewSPI(0, 0)

	if err != nil {
		log.Fatal(err)
	}

	// Convert the spi.Port into a spi.Conn so it can be used for communication.
	c, err := p.Connect(physic.MegaHertz, spi.Mode0, 8)

	if err != nil {
		log.Fatal(err)
	}
	return c
}

func delay(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

package icm20789

import (
	"log"
)

func (imu *imuIcm20789) setupGyro() {
	log.Println("SETUP IMU gyro")
	gyrosetup1, _ := imu.readByteFromRegister(GYRO_CONFIG)
	imu.writeRegister(GYRO_CONFIG, gyrosetup1|0b00110000)
	delay(1)
	gyrosetup2, _ := imu.readByteFromRegister(GYRO_CONFIG)
	log.Printf("GYRO_SETUP1: 0x%x, GYRO_SETUP2: 0x%x\n", gyrosetup1, gyrosetup2)
}

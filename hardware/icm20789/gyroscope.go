package icm20789

import (
	"log"
)

func (imu *imuICM20789) setupGyro() {
	log.Println("SETUP IMU gyro")
	gyrosetup1, _ := imu.readByteFromRegister(GYRO_CONFIG)
	imu.writeRegister(GYRO_CONFIG, gyrosetup1|0b00110000)
	delay(1)
	gyrosetup2, _ := imu.readByteFromRegister(GYRO_CONFIG)
	log.Printf("GYRO_SETUP1: 0x%x, GYRO_SETUP2: 0x%x\n", gyrosetup1, gyrosetup2)
}

func gyroscopeFullScale(fsr string) float64 {
	switch fsr {
	case "250dps":
		return GYRO_FULL_SCALE_250DPS
	case "500dps":
		return GYRO_FULL_SCALE_500DPS
	case "1000dps":
		return GYRO_FULL_SCALE_1000DPS
	case "2000dps":
		return GYRO_FULL_SCALE_2000DPS
	default:
		return GYRO_FULL_SCALE_250DPS
	}
}

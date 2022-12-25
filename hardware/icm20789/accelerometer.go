package icm20789

import "log"

// func accelerometerFullScale(fsr string) float64 {
// 	switch fsr {
// 	case "2g":
// 		return ACCEL_FULL_SCALE_2G
// 	case "4g":
// 		return ACCEL_FULL_SCALE_4G
// 	case "8g":
// 		return ACCEL_FULL_SCALE_8G
// 	case "16g":
// 		return ACCEL_FULL_SCALE_16G
// 	default:
// 		return ACCEL_FULL_SCALE_2G
// 	}
// }

func (imu *imuICM20789) setupAccelerometer(fullScaleMask byte) {
	log.Println("SETUP IMU Accelerometer")
	accelsetup1, _ := imu.readByteFromRegister(ACCEL_CONFIG)
	imu.writeRegister(ACCEL_CONFIG, accelsetup1|fullScaleMask)
	delay(1)
	accelsetup2, _ := imu.readByteFromRegister(ACCEL_CONFIG)
	log.Printf("ACCEL_CONFIG1: 0x%x, ACCEL_CONFIG2: 0x%x\n", accelsetup1, accelsetup2)
}

func accelerometerFullScale(fullScale string) (float64, byte) {
	fullScaleCoefficient := ACCEL_FULL_SCALE_2G
	fullScaleMask := byte(0b00000000)
	switch fullScale {
	case "2g":
		fullScaleCoefficient = ACCEL_FULL_SCALE_2G
		fullScaleMask = byte(0b00000000)
	case "4g":
		fullScaleCoefficient = ACCEL_FULL_SCALE_4G
		fullScaleMask = byte(0b00010000)
	case "8g":
		fullScaleCoefficient = ACCEL_FULL_SCALE_8G
		fullScaleMask = byte(0b00100000)
	case "16g":
		fullScaleCoefficient = ACCEL_FULL_SCALE_16G
		fullScaleMask = byte(0b00110000)
	}
	return fullScaleCoefficient, fullScaleMask
}

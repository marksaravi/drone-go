package icm20789

func (imu *imuICM20789) setupAccelerometer(fullScaleMask byte) {
	accelsetup1, _ := imu.readByteFromRegister(ACCEL_CONFIG)
	imu.writeRegister(ACCEL_CONFIG, accelsetup1|fullScaleMask)
	delay(1)
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

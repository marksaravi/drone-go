package icm20789

func (m *memsIcm20789) setupGyroscope(fullScaleMask byte) {
	gyrosetup1, _ := m.readByteFromRegister(ADRESS_GYRO_CONFIG)
	m.writeRegister(ADRESS_GYRO_CONFIG, gyrosetup1|fullScaleMask)
	delay(1)
}

func gyroscopeFullScale(fullScale string) (float64, byte) {
	fullScaleCoefficient := GYRO_FULL_SCALE_250DPS
	fullScaleMask := byte(0b00000000)
	switch fullScale {
	case "250dps":
		fullScaleCoefficient = GYRO_FULL_SCALE_250DPS
		fullScaleMask = byte(0b00000000)
	case "500dps":
		fullScaleCoefficient = GYRO_FULL_SCALE_500DPS
		fullScaleMask = byte(0b00010000)
	case "1000dps":
		fullScaleCoefficient = GYRO_FULL_SCALE_1000DPS
		fullScaleMask = byte(0b00100000)
	case "2000dps":
		fullScaleCoefficient = GYRO_FULL_SCALE_2000DPS
		fullScaleMask = byte(0b00110000)
	}
	return fullScaleCoefficient, fullScaleMask
}

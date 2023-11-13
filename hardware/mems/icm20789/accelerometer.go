package icm20789

const (
	ADDRESS_ACCEL_XOUT_H byte = 0x3B
	ADDRESS_ACCEL_CONFIG byte = 0x1C
	ADDRESS_XA_OFFSH     byte = 0x77
	ADDRESS_XA_OFFSL     byte = 0x78
	ADDRESS_YA_OFFSH     byte = 0x7A
	ADDRESS_YA_OFFSL     byte = 0x7B
	ADDRESS_ZA_OFFSH     byte = 0x7D
	ADDRESS_ZA_OFFSL     byte = 0x7E
)

const (
	ACCEL_CONFIG_DISABLE_SELF_TESTS  byte = 0b00000000
	ACCEL_CONFIG_2G                  byte = 0b00000000
	ACCEL_CONFIG_4G                  byte = 0b00001000
	ACCEL_CONFIG_8G                  byte = 0b00010000
	ACCEL_CONFIG_16G                 byte = 0b00011000
	ACCEL_CONFIG2_FIFO_SIZE_512      byte = 0b00000000
	ACCEL_CONFIG2_DEC2_CFG_4_SAMPLE  byte = 0b00000000
	ACCEL_CONFIG2_DEC2_CFG_8_SAMPLE  byte = 0b00010000
	ACCEL_CONFIG2_DEC2_CFG_16_SAMPLE byte = 0b00100000
	ACCEL_CONFIG2_DEC2_CFG_32_SAMPLE byte = 0b00110000
	ACCEL_CONFIG2_ACCEL_FCHOICE_B    byte = 0b00001000 //3-dB BW (Hz) 1046.0
	ACCEL_CONFIG2_A_DLPF_CFG_3dB_BW0 byte = 0b00000000 //3-dB BW (Hz) 218.1
	ACCEL_CONFIG2_A_DLPF_CFG_3dB_BW1 byte = 0b00000001 //3-dB BW (Hz) 218.1
	ACCEL_CONFIG2_A_DLPF_CFG_3dB_BW2 byte = 0b00000010 //3-dB BW (Hz) 99.0
	ACCEL_CONFIG2_A_DLPF_CFG_3dB_BW3 byte = 0b00000011 //3-dB BW (Hz) 44.8
	ACCEL_CONFIG2_A_DLPF_CFG_3dB_BW4 byte = 0b00000100 //3-dB BW (Hz) 21.2
	ACCEL_CONFIG2_A_DLPF_CFG_3dB_BW5 byte = 0b00000101 //3-dB BW (Hz) 10.2
	ACCEL_CONFIG2_A_DLPF_CFG_3dB_BW6 byte = 0b00000110 //3-dB BW (Hz) 5.1
	ACCEL_CONFIG2_A_DLPF_CFG_3dB_BW7 byte = 0b00000111 //3-dB BW (Hz) 420.0
)

const (
	ACCEL_FULL_SCALE_2G  float64 = 16384
	ACCEL_FULL_SCALE_4G  float64 = 8192
	ACCEL_FULL_SCALE_8G  float64 = 4096
	ACCEL_FULL_SCALE_16G float64 = 2048
)

func (m *memsIcm20789) setupAccelerometer(fullScaleMask byte) {
	accelsetup1, _ := m.readByteFromRegister(ADDRESS_ACCEL_CONFIG)
	m.writeRegister(ADDRESS_ACCEL_CONFIG, accelsetup1|fullScaleMask)
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

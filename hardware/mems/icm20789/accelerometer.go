package icm20789

const (
	ADDRESS_ACCEL_XOUT_H byte = 0x3B
	ADDRESS_ACCEL_XOUT_L byte = 0x3C
	ADDRESS_ACCEL_YOUT_H byte = 0x3D
	ADDRESS_ACCEL_YOUT_L byte = 0x3E
	ADDRESS_ACCEL_ZOUT_H byte = 0x3F
	ADDRESS_ACCEL_ZOUT_L byte = 0x40

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
	ACCEL_CONFIG2_FIFO_SIZE_512      byte = 0b00000000
	ACCEL_CONFIG2_ACCEL_FCHOICE_B    byte = 0b00001000 //3-dB BW (Hz) 1046.0
)


var ACCEL_CONFIG_G = map[int]byte {
	2  : 0b00000000, // 2G
	4  : 0b00001000, // 4G
	8  : 0b00010000, // 8G
	16 : 0b00011000, // 16G
}

var ACCEL_FULL_SCALE_G = map[int]float64 {
	2  : 16384,
	4  : 8192,
	8  : 4096,
	16 : 2048,
}

var DEC2_CFG_N_SAMPLE = map[int]byte {
	4  : 0b00000000, // 4  sample
	8  : 0b00010000, // 8  sample
	16 : 0b00100000, // 16 sample
	32 : 0b00110000, // 32 sample
}

var DLPF_CFG_3dB_BW = map[int]byte {
	0 : 0b00000000, //3-dB BW (Hz) 218.1
    1 : 0b00000001, //3-dB BW (Hz) 218.1
    2 : 0b00000010, //3-dB BW (Hz) 99.0
    3 : 0b00000011, //3-dB BW (Hz) 44.8
    4 : 0b00000100, //3-dB BW (Hz) 21.2
    5 : 0b00000101, //3-dB BW (Hz) 10.2
    6 : 0b00000110, //3-dB BW (Hz) 5.1
    7 : 0b00000111, //3-dB BW (Hz) 420.0
}

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

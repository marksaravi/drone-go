package icm20789

import (
	"fmt"
	"github.com/marksaravi/drone-go/hardware/mems"
)

const (
	ADDRESS_ACCEL_CONFIG  byte = 0x1C
	ADDRESS_ACCEL_CONFIG2 byte = 0x1D

	ADDRESS_ACCEL_XOUT_H  byte = 0x3B
	ADDRESS_ACCEL_XOUT_L  byte = 0x3C
	ADDRESS_ACCEL_YOUT_H  byte = 0x3D
	ADDRESS_ACCEL_YOUT_L  byte = 0x3E
	ADDRESS_ACCEL_ZOUT_H  byte = 0x3F
	ADDRESS_ACCEL_ZOUT_L  byte = 0x40
 
	ADDRESS_XA_OFFSH      byte = 0x77
	ADDRESS_XA_OFFSL      byte = 0x78
	ADDRESS_YA_OFFSH      byte = 0x7A
	ADDRESS_YA_OFFSL      byte = 0x7B
	ADDRESS_ZA_OFFSH      byte = 0x7D
	ADDRESS_ZA_OFFSL      byte = 0x7E
)

const (
	ACCEL_CONFIG_DISABLE_SELF_TESTS  byte = 0b00000000
)


var ACCEL_CONFIG_G = map[string]byte {
	"2g"  : 0b00000000, // 2G
	"4g"  : 0b00001000, // 4G
	"8g"  : 0b00010000, // 8G
	"16g" : 0b00011000, // 16G
}

var ACCEL_FULL_SCALE_G = map[string]float64 {
	"2g"  : 16384,
	"4g"  : 8192,
	"8g"  : 4096,
	"16g" : 2048,
}

var ACCEL_CONFIG2_FIFO_SIZE = map[int]byte {
	512  : 0b00000000, // 512 byte
	1024 : 0b01000000, // 1 kb
	2048 : 0b10000000, // 2 kb
	4096 : 0b11000000, // 4 kb
}

var ACCEL_CONFIG2_DEC2_CFG_N_SAMPLE = map[int]byte {
	4  : 0b00000000, // 4  sample
	8  : 0b00010000, // 8  sample
	16 : 0b00100000, // 16 sample
	32 : 0b00110000, // 32 sample
}

var ACCEL_CONFIG2_DLPF_CFG_3dB_BW = map[string]byte {
	"1046.0hz" : 0b00001000, //3-dB BW (Hz) 1046.0
	"523.0hz"  : 0b00000000, //3-dB BW (Hz) 523.0
    "218.1hz"  : 0b00000001, //3-dB BW (Hz) 218.1
    "99.0hz"   : 0b00000010, //3-dB BW (Hz) 99.0
    "44.8hz"   : 0b00000011, //3-dB BW (Hz) 44.8
    "21.2hz"   : 0b00000100, //3-dB BW (Hz) 21.2
    "10.2hz"   : 0b00000101, //3-dB BW (Hz) 10.2
    "5.1hz"    : 0b00000110, //3-dB BW (Hz) 5.1
    "420.0hz"  : 0b00000111, //3-dB BW (Hz) 420.0
}

func (m *memsIcm20789) setupAccelerometer(accelFullScale string, numberOfSamples int, fifoSize int, lowPassFilterFrequency string) {
	config1 := ACCEL_CONFIG_DISABLE_SELF_TESTS | ACCEL_CONFIG_G[accelFullScale]
	config2 := ACCEL_CONFIG2_FIFO_SIZE[fifoSize] | 
		ACCEL_CONFIG2_DEC2_CFG_N_SAMPLE[numberOfSamples] |
		ACCEL_CONFIG2_DLPF_CFG_3dB_BW[lowPassFilterFrequency]
	m.writeRegister(ADDRESS_ACCEL_CONFIG, config1)
	delay(100)
	m.writeRegister(ADDRESS_ACCEL_CONFIG2,
		config2)
	delay(100)
	accelconfig1, _ := m.readByteFromRegister(ADDRESS_ACCEL_CONFIG)
	accelconfig2, _ := m.readByteFromRegister(ADDRESS_ACCEL_CONFIG2)
	fmt.Printf("Accelerometer config1: %b, confog2: %b\n", accelconfig1, accelconfig2)
}

func (m *memsIcm20789) memsDataToAccelerometer(memsData []byte) mems.XYZ {
	x:=float64(towsComplementUint8ToInt16(memsData[0], memsData[1])) / m.accelFullScale
	y:= float64(towsComplementUint8ToInt16(memsData[2], memsData[3])) / m.accelFullScale
	z:= float64(towsComplementUint8ToInt16(memsData[4], memsData[5])) / m.accelFullScale
	fmt.Printf("%0.2f,  %0.2f,  %0.2f  ", x,y,z)
	accel := mems.XYZ{
		X: x,
		Y: y,
		Z: z,
	}
	return accel
}

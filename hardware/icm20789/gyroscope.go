package icm20789

import (
	"fmt"

	"github.com/marksaravi/drone-go/hardware/mems"
)

const (
	ADRESS_GYRO_CONFIG byte = 0x1B
	ADDRESS_XG_OFFSH   byte = 0x13
	ADDRESS_XG_OFFSL   byte = 0x14
	ADDRESS_YG_OFFSH   byte = 0x15
	ADDRESS_YG_OFFSL   byte = 0x16
	ADDRESS_ZG_OFFSH   byte = 0x17
	ADDRESS_ZG_OFFSL   byte = 0x18
)

var GYRO_CONFIG_DPS = map[string]byte{
	"250dps":  0b00000000,
	"500dps":  0b00001000,
	"1000dps": 0b00010000,
	"2000dps": 0b00011000,
}

var GYRO_FULL_SCALE_DPS = map[string]float64{
	"250dps":  131,
	"500dps":  65.5,
	"1000dps": 32.8,
	"2000dps": 16.4,
}

func (m *memsIcm20789) setupGyroscope(fullScale string, xOffset, yOffset, zOffset uint16) {
	m.writeRegister(ADRESS_GYRO_CONFIG, 0b00000000|GYRO_CONFIG_DPS[fullScale])
	delay(100)
	m.setupGyroscopeOffsets(xOffset, yOffset, zOffset)
	gyroconfig1, _ := m.readByteFromRegister(ADRESS_GYRO_CONFIG)
	fmt.Printf("Gyroscope config: %b\n", gyroconfig1)
}

func (m *memsIcm20789) memsDataToGyroscope(memsData []byte) mems.DXYZ {
	dx := float64(towsComplementUint8ToInt16(memsData[0], memsData[1])) / m.gyroFullScale
	dy := float64(towsComplementUint8ToInt16(memsData[2], memsData[3])) / m.gyroFullScale
	dz := float64(towsComplementUint8ToInt16(memsData[4], memsData[5])) / m.gyroFullScale
	return mems.DXYZ{
		DX: dx,
		DY: dy,
		DZ: dz,
	}
}

func (m *memsIcm20789) setupGyroscopeOffsets(xOffset, yOffset, zOffset uint16) {
	m.setupGyroscopeOffset(ADDRESS_XG_OFFSH, ADDRESS_XG_OFFSL, xOffset)
	m.setupGyroscopeOffset(ADDRESS_YG_OFFSH, ADDRESS_YG_OFFSL, yOffset)
	m.setupGyroscopeOffset(ADDRESS_ZG_OFFSH, ADDRESS_ZG_OFFSL, zOffset)
	delay(100)
}

func GyroOffsetToHL(offset uint16) (higherBits, lowerBits byte) {
	lowerBits = byte(offset & 0b0000000011111111)
	higherBits = byte(offset >> 8 & 0b0000000011111111)
	return
}

func GyroHLtoOffset(higherBits, lowerBits byte) uint16 {
	h := uint16(higherBits)
	l := uint16(lowerBits)
	h = h << 8
	return h | l
}

func (m *memsIcm20789) setupGyroscopeOffset(addressHigh, addressLow byte, offset uint16) {
	higherBits, lowerBits := GyroOffsetToHL(offset)
	m.writeRegister(addressLow, lowerBits)
	m.writeRegister(addressHigh, higherBits)

}

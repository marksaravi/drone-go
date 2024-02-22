package icm20789

import (
	"fmt"

	"github.com/marksaravi/drone-go/hardware/mems"
)

var GYRO_CONFIG_DPS = map[string]byte {
	"250dps"  : 0b00000000,
	"500dps"  : 0b00001000,
	"1000dps" : 0b00010000,
	"2000dps" : 0b00011000,
}

var GYRO_FULL_SCALE_DPS = map[string]float64 {
	"250dps"  : 131,
	"500dps"  : 65.5,
	"1000dps" : 32.8,
	"2000dps" : 16.4,
}

func (m *memsIcm20789) setupGyroscope(fullScale string) {
	m.writeRegister(ADRESS_GYRO_CONFIG, 0b00000000|GYRO_CONFIG_DPS[fullScale])
	delay(100)
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
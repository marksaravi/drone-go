package icm20789

import (
	"fmt"
	"log"
)

func (imu *imuIcm20789) Initialize(gyroFullScale string) {
	// imu.writeSPI(0x6B, []byte{0x01})
	// imu.writeSPI(0x6B, []byte{0x01})

	// delay(5)
	// imu.writeSPI(0x6A, []byte{0x10})
	// imu.writeSPI(0x6B, []byte{0x41})

	// delay(5)
	// imu.writeSPI(0x6C, []byte{0x3f})
	// imu.writeSPI(0xF5, []byte{0x00})
	// imu.writeSPI(0x19, []byte{0x09})
	// imu.writeSPI(0xEA, []byte{0x00})
	// imu.writeSPI(0x6B, []byte{0x01})

	// delay(5)
	// imu.writeSPI(0x6A, []byte{0x10})
	// imu.writeSPI(0x6B, []byte{0x41})

	// delay(5)
	// imu.writeSPI(0x6B, []byte{0x01})

	// delay(5)
	// imu.writeSPI(0x23, []byte{0x00})
	// imu.writeSPI(0x6B, []byte{0x41})

	// delay(5)
	// imu.writeSPI(0x1D, []byte{0xC0})
	// imu.writeSPI(0x6B, []byte{0x01})

	// delay(5)
	// imu.writeSPI(0x1A, []byte{0xC0})
	// imu.writeSPI(0x6B, []byte{0x41})

	// delay(5)
	// imu.writeSPI(0x38, []byte{0x01})

	// delay(5)
	// imu.writeSPI(0x6A, []byte{0x10})
	// imu.writeSPI(0x6B, []byte{0x41})

	// delay(5)
	// imu.writeSPI(0x6B, []byte{0x01})

	// delay(5)
	// imu.writeSPI(0x23, []byte{0x00})
	// imu.writeSPI(0x6B, []byte{0x41})

	// delay(5)
	// imu.writeSPI(0x6B, []byte{0x41})
	// imu.writeSPI(0x6C, []byte{0x3f})
	// imu.writeSPI(0x6B, []byte{0x41})

	// imu.writeSPI(0x6B, []byte{0x01})

	// delay(5)
	// imu.writeSPI(0x6A, []byte{0x10})
	// imu.writeSPI(0x6B, []byte{0x41})

	// delay(5)
	// imu.writeSPI(0x6B, []byte{0x01})

	// delay(5)
	// imu.writeSPI(0x23, []byte{0x00})
	// imu.writeSPI(0x6B, []byte{0x41})

	// delay(5)
	// // spi_dev->read_registers(0x6A, &v, 1);
	// // printf("reg 0x6A=0x%02x\n", v);
	// imu.writeSPI(0x6B, []byte{0x01})

	// delay(5)
	// imu.writeSPI(0x6A, []byte{0x10})
	// imu.writeSPI(0x6B, []byte{0x41})

	// delay(5)
	// imu.writeSPI(0x6B, []byte{0x01})

	// delay(5)
	// imu.writeSPI(0x23, []byte{0x00})
	// imu.writeSPI(0x6B, []byte{0x41})

	// delay(5)
	// imu.writeSPI(0x6B, []byte{0x41})
	// imu.writeSPI(0x6C, []byte{0x3f})
	// imu.writeSPI(0x6B, []byte{0x41})

	// imu.writeSPI(0x6B, []byte{0x01})

	// delay(5)
	// imu.writeSPI(0x6A, []byte{0x10})
	// imu.writeSPI(0x6B, []byte{0x41})

	// delay(5)
	// imu.writeSPI(0x6B, []byte{0x01})

	// delay(5)
	// imu.writeSPI(0x23, []byte{0x00})
	// imu.writeSPI(0x6B, []byte{0x41})

	imu.setAccelConfigs(gyroFullScale)
}

func (imu *imuIcm20789) setAccelConfigs(fullScale string) {
	const ADDRESS byte = 0x1B
	rwaconfig, _ := imu.readSPI(ADDRESS, 1)
	fmt.Println("RAW: ", rwaconfig[0])
	newValue := byte(0b01111111)
	writeErr := imu.writeSPI(ADDRESS, []byte{newValue})
	if writeErr != nil {
		log.Fatalf("can't write gyroscope config %v", writeErr)
	}
	delay(5)
	checkConfig, _ := imu.readSPI(ADDRESS, 1)
	fmt.Println(checkConfig[0], newValue)
}
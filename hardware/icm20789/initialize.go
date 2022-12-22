package icm20789

import (
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
	rwaconfig, firstReadErr := imu.readSPI(ACCEL_CONFIG, 1)
	newconfig := rwaconfig[0]
	if firstReadErr != nil {
		log.Fatalf("can't read accel config %v", firstReadErr)
	}
	switch fullScale {
	case "2g":
		newconfig = newconfig | ACCEL_CONFIG_MASK_FULL_SCALE_2G
	case "4g":
		newconfig = newconfig | ACCEL_CONFIG_MASK_FULL_SCALE_4G
	case "8g":
		newconfig = newconfig | ACCEL_CONFIG_MASK_FULL_SCALE_8G
	case "16":
		newconfig = newconfig | ACCEL_CONFIG_MASK_FULL_SCALE_16G
	default:
		log.Printf("incorrect accel config, using default 250dps")
	}
	log.Printf("new accel config is: %d\n", newconfig)

	writeErr := imu.writeSPI(GYRO_CONFIG, []byte{newconfig, 0})
	if writeErr != nil {
		log.Fatalf("can't write gyroscope config %v", writeErr)
	}
	delay(5)
	checkConfig, checkErr := imu.readSPI(GYRO_CONFIG, 1)
	if checkErr != nil {
		log.Fatalf("can't read gyroscope config %v", checkErr)
	}
	if checkConfig[0] != newconfig {
		log.Fatalf("can't write gyroscope config %d!=%d", checkConfig[0], newconfig)
	}
	log.Printf("successful accel configuration: %b\n", checkConfig[0])
}

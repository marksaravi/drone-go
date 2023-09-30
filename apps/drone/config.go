package drone

import (
	"log"
	"os"

	"encoding/json"

	"github.com/marksaravi/drone-go/hardware"
)

type InertialDeviceConfigs struct {
	FullScale string `json:"full_scale"`
	Offsets   struct {
		X uint16 `json:"x"`
		Y uint16 `json:"y"`
		Z uint16 `json:"z"`
	} `json:"offsets"`
}

type DroneConfigs struct {
	IMU struct {
		DataPerSecond     int     `json:"data-per-second"`
		FilterCoefficient float64 `json:"filter-coefficient"`
		MEMS              struct {
			Accelerometer InertialDeviceConfigs   `json:"accelerometer"`
			Gyroscope     InertialDeviceConfigs   `json:"gyroscope"`
			SPI           hardware.SPIConnConfigs `json:"spi"`
		} `json:"icm20789"`
	} `json:"imu"`
	RemoteControl struct {
		CommandsPerSecond int `json:"commands-per-second"`
		Radio             struct {
			RxTxAddress string                  `json:"rx-tx-address"`
			SPI         hardware.SPIConnConfigs `json:"spi"`
		} `json:"radio"`
	} `json:"remote-control"`
}

func ReadConfigs(configPath string) DroneConfigs {
	content, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	var configs DroneConfigs
	json.Unmarshal([]byte(content), &configs)
	return configs
}

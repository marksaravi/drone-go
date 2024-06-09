package drone

import (
	"log"
	"os"

	"encoding/json"

	"github.com/marksaravi/drone-go/devices/imu"
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

type PIDConfigs struct {
	P                   float64 `json:"p"`
	I                   float64 `json:"i"`
	D                   float64 `json:"d"`
	MaxRotationError    float64 `json:"max-rot-error"`
	MaxIntegrationValue float64 `json:"max-i-value"`
	MaxWeightedSum      float64 `json:"max-weighted-sum"`
	CalibrationMode     bool    `json:"calibration-mode"`
	CalibrationIncP     float64 `json:"calibration-p-inc"`
	CalibrationIncI     float64 `json:"calibration-i-inc"`
	CalibrationIncD     float64 `json:"calibration-d-inc"`
}

type DroneConfigs struct {
	PID PIDConfigs  `json:"pid"`
	IMU imu.Configs `json:"imu"`
	ESC struct {
		DataPerSecond int `json:"data-per-second"`
	} `json:"esc"`
	Commands struct {
		RollMidValue  int     `json:"roll-mid-value"`
		PitchMidValue int     `json:"pitch-mid-value"`
		YawMidValue   int     `json:"yaw-mid-value"`
		RotationRange float64 `json:"rotation-range-deg"`
		MaxThrottle   float64 `json:"max-throttle"`
	} `json:"commands"`
	RemoteControl struct {
		CommandsPerSecond int `json:"commands-per-second"`
		Radio             struct {
			RxTxAddress string                  `json:"rx-tx-address"`
			SPI         hardware.SPIConnConfigs `json:"spi"`
		} `json:"radio"`
	} `json:"remote-control"`
	Plotter struct {
		Active bool `json:"active"`
	} `json:"plotter"`
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

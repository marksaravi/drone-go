package drone

import (
	"log"
	"os"

	"encoding/json"

	"github.com/marksaravi/drone-go/devices/imu"
	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/pid"
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
	PID struct {
		ARM_0_2 pid.PIDConfigs `json:"arm_0_2"`
		ARM_1_3 pid.PIDConfigs `json:"arm_1_3"`
		Yaw     pid.PIDConfigs `json:"yaw"`
	} `json:"pid"`
	IMU imu.Configs `json:"imu"`
	ESC struct {
		DataPerSecond     int     `json:"data-per-second"`
		MaxOutputThrottle float64 `json:"max-output-throttle"`
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

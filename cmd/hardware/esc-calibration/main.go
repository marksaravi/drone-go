package main

import (
	"github.com/MarkSaravi/drone-go/config"
	"github.com/MarkSaravi/drone-go/hardware/pca9685"
)

func main() {
	appConfig := config.ReadConfigs()
	pca9685.Calibrate(appConfig.Hardware.PCA9685)
}

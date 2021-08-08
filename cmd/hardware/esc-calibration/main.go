package main

import (
	"github.com/MarkSaravi/drone-go/hardware/pca9685"
	"github.com/MarkSaravi/drone-go/utils"
)

func main() {
	appConfig := utils.ReadConfigs()
	pca9685.Calibrate(appConfig.Hardware.PCA9685)
}

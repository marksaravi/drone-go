package main

import (
	"github.com/MarkSaravi/drone-go/modules/udplogger"
	"github.com/MarkSaravi/drone-go/types"
)

func initUdpLogger(appConfig ApplicationConfig) types.UdpLogger {
	udpLogger := udplogger.CreateUdpLogger(appConfig.UDP, appConfig.Flight.Imu.ImuDataPerSecond)
	return udpLogger
}

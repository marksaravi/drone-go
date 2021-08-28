package flightcontrol

import (
	"fmt"

	"github.com/MarkSaravi/drone-go/models"
)

type imu interface {
	Read() (models.ImuRotations, bool)
}

type udpLogger interface {
	Send(models.ImuRotations)
}

type flightControl struct {
	imu       imu
	udpLogger udpLogger
}

func NewFlightControl(imu imu, udpLogger udpLogger) *flightControl {
	return &flightControl{
		imu:       imu,
		udpLogger: udpLogger,
	}
}

func (fc *flightControl) Start() {
	fmt.Println("Starting Flight Control")
	for {
		data, canRead := fc.imu.Read()
		if canRead {
			fc.udpLogger.Send(data)
		}
	}
}

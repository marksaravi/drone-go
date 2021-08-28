package flightcontrol

import (
	"fmt"

	"github.com/MarkSaravi/drone-go/models"
)

type imu interface {
	Read() (models.ImuRotations, bool)
}

type flightControl struct {
	imu imu
}

func NewFlightControl(imu imu) *flightControl {
	return &flightControl{
		imu: imu,
	}
}

func (fc *flightControl) Start() {
	fmt.Println("Starting Flight Control")
	for {
		fc.imu.Read()
	}
}

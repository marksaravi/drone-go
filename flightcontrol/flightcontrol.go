package flightcontrol

import (
	"fmt"
)

type imu interface {
	Read() (canRead bool)
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

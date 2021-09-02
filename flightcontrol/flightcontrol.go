package flightcontrol

import (
	"fmt"

	"github.com/MarkSaravi/drone-go/models"
)

type radio interface {
	IsDataAvailable() bool
	ReceiverOn()
	ReceiveFlightData() models.FlightData
	TransmitterOn()
	TransmitFlightData(models.FlightData) error
}

type imu interface {
	Read() (models.ImuRotations, bool)
}

type udpLogger interface {
	Send(models.ImuRotations)
}

type flightControl struct {
	imu       imu
	radio     radio
	udpLogger udpLogger
}

func NewFlightControl(imu imu, radio radio, udpLogger udpLogger) *flightControl {
	return &flightControl{
		imu:       imu,
		radio:     radio,
		udpLogger: udpLogger,
	}
}

func (fc *flightControl) Start() {
	fmt.Println("Starting Flight Control")
	fc.radio.ReceiverOn()
	for {
		data, canRead := fc.imu.Read()
		if canRead {
			fc.udpLogger.Send(data)
		}
		if fc.radio.IsDataAvailable() {
			fmt.Println(fc.radio.ReceiveFlightData())
		}
	}
}

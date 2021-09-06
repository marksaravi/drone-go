package flightcontrol

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/models"
)

type radio interface {
	ReceiverOn()
	ReceiveFlightData() (models.FlightData, bool)
	TransmitterOn()
	TransmitFlightData(models.FlightData) error
}

type imu interface {
	ReadRotations() models.ImuRotations
	ResetTime()
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
	var dataPerSecond int = 3200
	readingInterval := time.Second / time.Duration(dataPerSecond)
	// fc.radio.ReceiverOn()
	fc.imu.ResetTime()
	lastReadingTime := time.Now()
	for {
		now := time.Now()
		if now.Sub(lastReadingTime) >= readingInterval {
			rotations := fc.imu.ReadRotations()
			fc.udpLogger.Send(rotations)
			printFlightData(rotations)
			lastReadingTime = now
		}
		// flightData, isAvalable := fc.radio.ReceiveFlightData()
		// if isAvalable {
		// 	printFlightData(flightData)
		// 	fc.radio.TransmitterOn()
		// 	fc.radio.TransmitFlightData(models.FlightData{
		// 		Id:              flightData.Id,
		// 		IsRemoteControl: false,
		// 		IsDrone:         true,
		// 	})
		// 	fc.radio.ReceiverOn()
		// }
	}
}

var lastPrint time.Time = time.Now()

func printFlightData(flightData models.ImuRotations) {
	if time.Since(lastPrint) < time.Second/4 {
		return
	}
	lastPrint = time.Now()
	fmt.Println(flightData)
}

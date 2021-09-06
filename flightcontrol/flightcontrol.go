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
	imuDataPerSecond int
	imu              imu
	radio            radio
	udpLogger        udpLogger
}

func NewFlightControl(imuDataPerSecond int, imu imu, radio radio, udpLogger udpLogger) *flightControl {
	return &flightControl{
		imuDataPerSecond: imuDataPerSecond,
		imu:              imu,
		radio:            radio,
		udpLogger:        udpLogger,
	}
}

func (fc *flightControl) Start() {
	fmt.Println("Starting Flight Control")
	readingInterval := time.Second / time.Duration(fc.imuDataPerSecond)
	fc.radio.ReceiverOn()
	commandChannel := NewCommandChannel(fc.radio)
	fc.imu.ResetTime()
	lastReadingTime := time.Now()
	var flightCommand models.FlightData

	var counter int = 0
	sampleStart := time.Now()
	for {
		now := time.Now()
		select {
		case flightCommand = <-commandChannel:
			acknowledge(flightCommand, fc.radio)
		default:
			if now.Sub(lastReadingTime) >= readingInterval {
				rotations := fc.imu.ReadRotations()
				fc.udpLogger.Send(rotations)
				lastReadingTime = now
				counter++
				if counter == fc.imuDataPerSecond {
					fmt.Println(time.Since(sampleStart))
					sampleStart = time.Now()
					counter = 0
				}
			}
		}
	}
}

func NewCommandChannel(r radio) chan models.FlightData {
	radioChannel := make(chan models.FlightData, 10)
	go func(r radio, c chan models.FlightData) {
		ticker := time.NewTicker(time.Second / 40)
		for range ticker.C {
			if d, isOk := r.ReceiveFlightData(); isOk {
				c <- d
			}
		}
	}(r, radioChannel)
	return radioChannel
}

var lastCommandPrinted = time.Now()

func acknowledge(fd models.FlightData, radio radio) {
	if time.Since(lastCommandPrinted) >= time.Second {
		lastCommandPrinted = time.Now()
		fmt.Println(fd)
	}
	radio.TransmitterOn()
	radio.TransmitFlightData(fd)
	radio.ReceiverOn()
}

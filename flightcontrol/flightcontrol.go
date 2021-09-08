package flightcontrol

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/models"
	"github.com/MarkSaravi/drone-go/utils"
)

const (
	commandTimeCorrectionPercent float32 = 2.5
	escTimeCorrectionPercent     float32 = 3.5
	imuTimeCorrectionPercent     float32 = 8.5
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
type esc interface {
	Off()
	On()
	SetThrottles(map[uint8]float32)
}

type udpLogger interface {
	Send(models.ImuRotations)
}

type flightControl struct {
	imuDataPerSecond   int
	escUpdatePerSecond int
	imu                imu
	esc                esc
	radio              radio
	udpLogger          udpLogger
}

func NewFlightControl(imuDataPerSecond int, escUpdatePerSecond int, imu imu, esc esc, radio radio, udpLogger udpLogger) *flightControl {
	return &flightControl{
		imuDataPerSecond:   imuDataPerSecond,
		escUpdatePerSecond: escUpdatePerSecond,
		imu:                imu,
		esc:                esc,
		radio:              radio,
		udpLogger:          udpLogger,
	}
}

func (fc *flightControl) Start() {
	fmt.Printf("IMU data/s: %d, ESC refresh/s: %d\n", fc.imuDataPerSecond, fc.escUpdatePerSecond)
	fc.radio.ReceiverOn()
	imuDataChannel := newImuDataChannel(fc.imu, fc.imuDataPerSecond)
	escThrottleControlChannel := newEscThrottleControlChannel(fc.esc)
	escRefresher := utils.NewTicker("esc", fc.escUpdatePerSecond, true, escTimeCorrectionPercent)
	commandChannel := newCommandChannel(fc.radio)
	fc.esc.On()
	defer fc.esc.Off()
	time.Sleep(3 * time.Second)
	var throttle float32 = 0
	var running bool = true
	var rotations models.ImuRotations
	for running {
		select {
		case fd := <-commandChannel:
			throttle = fd.Throttle / 5 * 10
			if fd.IsMotorsEngaged {
				running = false
			}
		case rotations = <-imuDataChannel:
			fc.udpLogger.Send(rotations)
		case <-escRefresher:
			escThrottleControlChannel <- map[uint8]float32{
				0: throttle,
				1: throttle,
				2: throttle,
				3: throttle,
			}
		default:
			utils.Idle()
		}
	}
}

func newEscThrottleControlChannel(escdevice esc) chan map[uint8]float32 {
	escChannel := make(chan map[uint8]float32, 10)
	go func(escdev esc, ch chan map[uint8]float32) {
		var throttles map[uint8]float32
		for {
			select {
			case throttles = <-ch:
				escdev.SetThrottles(throttles)
			default:
				utils.Idle()
			}
		}
	}(escdevice, escChannel)
	return escChannel
}

func newImuDataChannel(imudev imu, dataPerSecond int) <-chan models.ImuRotations {
	imuDataChannel := make(chan models.ImuRotations, 10)
	go func(imudev imu, ch chan models.ImuRotations) {
		ticker := utils.NewTicker("imu", dataPerSecond, true, imuTimeCorrectionPercent)
		for range ticker {
			ch <- imudev.ReadRotations()
		}
	}(imudev, imuDataChannel)
	return imuDataChannel
}

func newCommandChannel(r radio) <-chan models.FlightData {
	radioChannel := make(chan models.FlightData, 10)
	go func(r radio, c chan models.FlightData) {
		ticker := utils.NewTicker("command", 40, true, commandTimeCorrectionPercent)
		for range ticker {
			if d, isOk := r.ReceiveFlightData(); isOk {
				c <- d
			}
			utils.Idle()
		}
	}(r, radioChannel)
	go func(c chan models.FlightData) {
		fmt.Scanln()
		c <- models.FlightData{
			IsMotorsEngaged: true,
		}
	}(radioChannel)

	return radioChannel
}

func acknowledge(fd models.FlightData, radio radio) {
	radio.TransmitterOn()
	radio.TransmitFlightData(fd)
	radio.ReceiverOn()
}

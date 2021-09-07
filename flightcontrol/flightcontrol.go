package flightcontrol

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/models"
	"github.com/MarkSaravi/drone-go/utils"
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
	imuDataChannel := newImuDataChannel(fc.imu, fc.imuDataPerSecond)
	escThrottleControlChannel := newEscThrottleControlChannel(fc.esc)
	escRefresher := utils.NewTicker(fc.escUpdatePerSecond)
	imustart := time.Now()
	imucounter := 0
	escstart := time.Now()
	esccounter := 0
	for {
		select {
		case <-imuDataChannel:
			imucounter++
			if imucounter == fc.imuDataPerSecond {
				fmt.Println("imu: ", time.Since(imustart))
				imustart = time.Now()
				imucounter = 0
			}
		case <-escRefresher:
			esccounter++
			if esccounter == fc.escUpdatePerSecond {
				fmt.Println("esc: ", time.Since(escstart))
				escstart = time.Now()
				esccounter = 0
			}
			escThrottleControlChannel <- 13425598
		default:
		}
	}
}

func newEscThrottleControlChannel(escdevice esc) chan int64 {
	escChannel := make(chan int64, 10)
	go func(escdev esc, ch chan int64) {
		var throttles int64
		start := time.Now()
		for {
			select {
			case throttles = <-ch:
				var motor uint8
				for motor = 0; motor < 4; motor++ {
					// escdev.SetThrottle(int(motor), throttles[motor])
				}
				if time.Since(start) >= time.Second {
					fmt.Println(throttles)
					start = time.Now()
				}
			default:
			}
			time.Sleep(time.Nanosecond)
		}
	}(escdevice, escChannel)
	return escChannel
}

func newImuDataChannel(imudev imu, dataPerSecond int) <-chan models.ImuRotations {
	imuDataChannel := make(chan models.ImuRotations, 10)
	go func(imudev imu, ch chan models.ImuRotations) {
		ticker := utils.NewTicker(dataPerSecond)
		for range ticker {
			ch <- imudev.ReadRotations()
		}
	}(imudev, imuDataChannel)
	return imuDataChannel
}

func newCommandChannel(r radio) <-chan models.FlightData {
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

func acknowledge(fd models.FlightData, radio radio) {
	radio.TransmitterOn()
	radio.TransmitFlightData(fd)
	radio.ReceiverOn()
}

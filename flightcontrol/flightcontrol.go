package flightcontrol

import (
	"fmt"
	"sync"
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
type esc interface {
	Off()
	On()
	SetThrottle(int, float32)
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
	readingInterval := time.Second / time.Duration(fc.imuDataPerSecond)
	imuUpdatePerEscUpdate := fc.imuDataPerSecond / fc.escUpdatePerSecond
	fmt.Printf("Starting Flight Control, imu dpr: %d, esc upeu: %d\n", fc.imuDataPerSecond, imuUpdatePerEscUpdate)
	fc.radio.ReceiverOn()
	commandChannel := NewCommandChannel(fc.radio)
	lastReadingTime := time.Now()
	var flightCommand models.FlightData
	var imuCheckCounter int = 0
	imuCheckTimer := time.Now()
	var escCheckCounter int = 0
	escCheckTimer := time.Now()
	var escUpdateCounter int = 0
	var escMutex sync.Mutex
	fc.imu.ResetTime()
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
				escUpdateCounter++
				imuCheckCounter++
				if escUpdateCounter == imuUpdatePerEscUpdate {
					escUpdateCounter = 0
					go func(el *sync.Mutex) {
						el.Lock()
						escCheckCounter++
						for i := 0; i < 4; i++ {
							fc.esc.SetThrottle(i, 0)
						}
						el.Unlock()
					}(&escMutex)
				}
				if imuCheckCounter == fc.imuDataPerSecond {
					fmt.Println("imu: ", time.Since(imuCheckTimer))
					imuCheckTimer = time.Now()
					imuCheckCounter = 0
				}
				if escCheckCounter == fc.escUpdatePerSecond {
					fmt.Println("esc: ", time.Since(escCheckTimer))
					escCheckTimer = time.Now()
					escCheckCounter = 0
				}
			}
		}
		time.Sleep(time.Microsecond)
	}
}

func NewCommandChannel(r radio) <-chan models.FlightData {
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

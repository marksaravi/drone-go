package models

import (
	"time"

	"github.com/marksaravi/drone-go/constants"
)

type XYZ struct {
	X float64
	Y float64
	Z float64
}

type Rotations struct {
	Roll, Pitch, Yaw float64
}

type ImuRotations struct {
	Accelerometer Rotations
	Gyroscope     Rotations
	Rotations     Rotations
	ReadTime      time.Time
	ReadInterval  time.Duration
}

type FlightCommands struct {
	PayloadType       uint8
	Roll              uint16
	Pitch             uint16
	Yaw               uint16
	Throttle          uint16
	ButtonFrontLeft   bool
	ButtonFrontRight  bool
	ButtonTopLeft     bool
	ButtonTopRight    bool
	ButtonBottomLeft  bool
	ButtonBottomRight bool
}

type Payload = [constants.RADIO_PAYLOAD_SIZE]byte

type RadioLink interface {
	Receive() (Payload, bool)
	Transmit(Payload) error
}

type Radio interface {
	Transmit(FlightCommands)
	GetReceiver() <-chan FlightCommands
	GetConnection() <-chan int
	CloseTransmitter()
}

type Logger interface {
	Send(ImuRotations)
	Close()
}

type Throttles = map[uint8]float32

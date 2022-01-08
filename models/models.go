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
	Accelerometer        Rotations
	Gyroscope            Rotations
	Rotations            Rotations
	TransformedRotations Rotations
	ReadTime             time.Time
	ReadInterval         time.Duration
}

type FlightCommandType = uint8

type FlightCommands struct {
	Type              FlightCommandType
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

type Radio interface {
	Transmit(FlightCommands)
	GetReceiver() <-chan FlightCommands
	GetConnection() <-chan int
	Close()
	SuppressLostConnection()
}

type Logger interface {
	Send(ImuRotations)
	Close()
}

type Throttles struct {
	Throttle         float64
	ControlVariables map[int]float64
}

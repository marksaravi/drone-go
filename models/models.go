package models

import "time"

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
	Id                uint32
	Time              int64
	Roll              float32
	Pitch             float32
	Yaw               float32
	Throttle          float32
	ButtonFrontLeft   bool
	ButtonFrontRight  bool
	ButtonTopLeft     bool
	ButtonTopRight    bool
	ButtonBottomLeft  bool
	ButtonBottomRight bool
}

type RadioLink interface {
	ReceiverOn()
	Receive() ([32]byte, bool)
	TransmitterOn()
	Transmit([32]byte) error
}

type Radio interface {
	Transmit(FlightCommands) bool
	GetReceiver() <-chan FlightCommands
	GetConnection() <-chan bool
}

type Throttles = map[uint8]float32

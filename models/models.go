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
	PayloadType       byte
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

type Payload = [32]byte

type RadioLink interface {
	ReceiverOn()
	Receive() (Payload, bool)
	TransmitterOn()
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

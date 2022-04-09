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

type RotationsAroundImuAxis struct {
	X, Y, Z float64
}

type ImuRotations struct {
	Accelerometer RotationsAroundImuAxis
	Gyroscope     RotationsAroundImuAxis
	Rotations     RotationsAroundImuAxis
	ReadTime      time.Time
	ReadInterval  time.Duration
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

type Logger interface {
	Send(ImuRotations)
	Close()
}

type Throttles struct {
	BaseThrottle float64
	Throttles    map[int]float64
}

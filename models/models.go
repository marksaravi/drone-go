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

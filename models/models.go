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

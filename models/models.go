package models

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
	ReadTime      int64
	ReadInterval  int64
}

type FlightCommands struct {
	Id               uint32
	Roll             float32
	Pitch            float32
	Yaw              float32
	Throttle         float32
	ButtonFrontLeft  bool
	ButtonFrontRight bool
	ButtonTopLeft    bool
	ButtonTopRight   bool
	ButtonDownLeft   bool
	ButtonDownRight  bool
}

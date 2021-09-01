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

type FlightData struct {
	Id       uint32
	Roll     float32
	Pitch    float32
	Yaw      float32
	Throttle float32
	Altitude float32
}

type JoystickData struct {
	Value     float32
	IsChanged bool
}

type ButtonData struct {
	Value     bool
	IsChanged bool
}

type RemoteControlData struct {
	Roll            JoystickData
	Pitch           JoystickData
	Yaw             JoystickData
	Throttle        JoystickData
	ButtonFrontLeft ButtonData
}

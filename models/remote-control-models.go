package models

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

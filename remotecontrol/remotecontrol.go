package remotecontrol

import "github.com/MarkSaravi/drone-go/modules/adcconverter"

type MotorsState int

const (
	Off MotorsState = iota
	On
	EmergencyOff
)

type RemoteData struct {
	Throttle    float32
	X           float32
	Y           float32
	Z           float32
	MotorsState MotorsState
}

type RemoteControl interface {
	ReadInputs() RemoteData
}

type remoteControl struct {
	adc adcconverter.AnalogToDigitalConverter
}

func NewRemoteControl(adc adcconverter.AnalogToDigitalConverter) RemoteControl {
	return &remoteControl{
		adc: adc,
	}
}

func (rc *remoteControl) ReadInputs() RemoteData {
	rc.adc.ReadInputVoltage(0)
	return RemoteData{
		Throttle:    0,
		X:           0,
		Y:           0,
		Z:           0,
		MotorsState: Off,
	}
}

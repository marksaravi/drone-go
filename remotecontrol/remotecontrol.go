package remotecontrol

import (
	"github.com/MarkSaravi/drone-go/hardware/mcp3008"
	"github.com/MarkSaravi/drone-go/modules/adcconverter"
)

type MotorsState int

const (
	Off MotorsState = iota
	On
	EmergencyOff
)

type RemoteControlConfig struct {
	MCP3008             mcp3008.MCP3008Config `yaml:"mcp3008"`
	XChannel            int                   `yaml:"x-channel"`
	YChannel            int                   `yaml:"y-channel"`
	ZChannel            int                   `yaml:"z-channel"`
	ThrottleChannel     int                   `yaml:"throttle-channel"`
	ReadyLight          string                `yaml:"ready-light-gpio"`
	EmergencyStopLight  string                `yaml:"emergency-stop-light-gpio"`
	EmergencyStopButton string                `yaml:"emergency-stop-button-gpio"`
}

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

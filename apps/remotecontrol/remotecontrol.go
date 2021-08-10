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
	VRef                float32               `yaml:"v-ref"`
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
	adc             adcconverter.AnalogToDigitalConverter
	vRef            float32
	xChannel        int
	yChannel        int
	zChannel        int
	throttleChannel int
}

func NewRemoteControl(adc adcconverter.AnalogToDigitalConverter, config RemoteControlConfig) RemoteControl {
	return &remoteControl{
		adc:             adc,
		vRef:            config.VRef,
		xChannel:        config.XChannel,
		yChannel:        config.YChannel,
		zChannel:        config.ZChannel,
		throttleChannel: config.ThrottleChannel,
	}
}

func (rc *remoteControl) ReadInputs() RemoteData {
	x, _ := rc.adc.ReadInputVoltage(rc.xChannel, rc.vRef)
	y, _ := rc.adc.ReadInputVoltage(rc.yChannel, rc.vRef)
	z, _ := rc.adc.ReadInputVoltage(rc.zChannel, rc.vRef)
	throttle, _ := rc.adc.ReadInputVoltage(rc.throttleChannel, rc.vRef)
	return RemoteData{
		Throttle:    throttle,
		X:           x,
		Y:           y,
		Z:           z,
		MotorsState: Off,
	}
}

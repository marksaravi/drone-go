package remotecontrol

import (
	"fmt"

	"github.com/MarkSaravi/drone-go/hardware/mcp3008"
	"github.com/MarkSaravi/drone-go/modules/adcconverter"
	"periph.io/x/periph/conn/gpio"
)

type MotorsState int

const (
	Off MotorsState = iota
	On
	EmergencyOff
)

type RemoteControlConfig struct {
	MCP3008          mcp3008.MCP3008Config `yaml:"mcp3008"`
	XChannel         int                   `yaml:"x-channel"`
	YChannel         int                   `yaml:"y-channel"`
	ZChannel         int                   `yaml:"z-channel"`
	ThrottleChannel  int                   `yaml:"throttle-channel"`
	ButtonFrontLeft  string                `yaml:"button-front-left"`
	ButtonFrontRight string                `yaml:"button-front-right"`
	ButtonTopLeft    string                `yaml:"button-top-left"`
	ButtonTopRight   string                `yaml:"button-top-right"`
	VRef             float32               `yaml:"v-ref"`
}

type RemoteData struct {
	Throttle float32
	X        float32
	Y        float32
	Z        float32
	Stop     bool
	Front    bool
}

type RemoteControl interface {
	ReadInputs() RemoteData
}

type remoteControl struct {
	adc              adcconverter.AnalogToDigitalConverter
	vRef             float32
	xChannel         int
	yChannel         int
	zChannel         int
	throttleChannel  int
	buttonFrontLeft  gpio.PinIn
	buttonFrontRight gpio.PinIn
	buttonTopLeft    gpio.PinIn
	buttonTopRight   gpio.PinIn
}

func NewRemoteControl(
	adc adcconverter.AnalogToDigitalConverter,
	buttonFrontLeft gpio.PinIn,
	buttonFrontRight gpio.PinIn,
	buttonTopLeft gpio.PinIn,
	buttonTopRight gpio.PinIn,
	config RemoteControlConfig,
) RemoteControl {
	return &remoteControl{
		adc:              adc,
		vRef:             config.VRef,
		xChannel:         config.XChannel,
		yChannel:         config.YChannel,
		zChannel:         config.ZChannel,
		throttleChannel:  config.ThrottleChannel,
		buttonFrontLeft:  buttonFrontLeft,
		buttonFrontRight: buttonFrontRight,
		buttonTopLeft:    buttonTopLeft,
		buttonTopRight:   buttonTopRight,
	}
}

func (rc *remoteControl) ReadInputs() RemoteData {
	x, _ := rc.adc.ReadInputVoltage(rc.xChannel, rc.vRef)
	y, _ := rc.adc.ReadInputVoltage(rc.yChannel, rc.vRef)
	z, _ := rc.adc.ReadInputVoltage(rc.zChannel, rc.vRef)
	throttle, _ := rc.adc.ReadInputVoltage(rc.throttleChannel, rc.vRef)
	topLeft := rc.buttonTopLeft.Read() == gpio.Low
	topRight := rc.buttonTopRight.Read() == gpio.Low
	frontLeft := rc.buttonFrontLeft.Read() == gpio.Low
	frontRight := rc.buttonFrontRight.Read() == gpio.Low
	fmt.Println(topLeft, topRight, frontLeft, frontRight)
	return RemoteData{
		X:        x,
		Y:        y,
		Z:        z,
		Throttle: throttle,
		Stop:     topLeft || topRight,
		Front:    frontLeft || frontRight,
	}
}

package devices

import "github.com/MarkSaravi/drone-go/models"

type analogtodigital interface {
	Read() float32
}

type joystickInput struct {
	input analogtodigital
	data  models.JoystickData
}

func (js *joystickInput) Read() models.JoystickData {
	pv := js.data.Value
	js.data.Value = js.input.Read()
	js.data.IsChanged = js.data.Value != pv
	return js.data
}

func NewJoystick(ad analogtodigital) *joystickInput {
	return &joystickInput{
		input: ad,
	}
}

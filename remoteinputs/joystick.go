package remoteinputs

import "github.com/MarkSaravi/drone-go/models"

type joystick interface {
	Read() float32
}

type joystickInput struct {
	input joystick
	data  models.JoystickData
}

func (js *joystickInput) read() {
	pv := js.data.Value
	js.data.Value = js.input.Read()
	js.data.IsChanged = js.data.Value != pv
}

func (ri *remoteInputs) readJoysticks() {
	ri.roll.read()
	ri.pitch.read()
	ri.yaw.read()
}

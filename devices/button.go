package devices

import "github.com/MarkSaravi/drone-go/models"

type gpioswitch interface {
	Read() bool
}

type button struct {
	input gpioswitch
	data  models.ButtonData
}

func (btn *button) Read() models.ButtonData {
	pv := btn.data.Value
	btn.data.Value = btn.input.Read()
	btn.data.IsChanged = btn.data.Value != pv
	return btn.data
}

func NewButton(input gpioswitch) *button {
	return &button{
		input: input,
	}
}

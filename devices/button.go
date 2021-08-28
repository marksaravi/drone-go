package devcices

import "github.com/MarkSaravi/drone-go/models"

type input interface {
	Read() bool
}

type button struct {
	input input
	data  models.ButtonData
}

func (btn *button) Read() models.ButtonData {
	pv := btn.data.Value
	btn.data.Value = btn.input.Read()
	btn.data.IsChanged = btn.data.Value != pv
	return btn.data
}

func NewButton(input input) *button {
	return &button{
		input: input,
	}
}

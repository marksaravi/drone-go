package remoteinputs

import "github.com/MarkSaravi/drone-go/models"

type button interface {
	Read() bool
}

type buttonInput struct {
	input button
	data  models.ButtonData
}

func (btn *buttonInput) read() {
	pv := btn.data.Value
	btn.data.Value = btn.input.Read()
	btn.data.IsChanged = btn.data.Value != pv
}

func (ri *remoteInputs) readStopButtons() {
	ri.buttonFrontLeft.read()
}

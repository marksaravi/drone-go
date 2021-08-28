package remoteinputs

type button interface {
	Read() bool
}

type buttonInput struct {
	input     button
	value     bool
	isChanged bool
}

func (btn *buttonInput) read() {
	pv := btn.value
	btn.value = btn.input.Read()
	btn.isChanged = btn.value != pv
}

func (ri *remoteInputs) readStopButtons() {
	ri.buttonFrontLeft.read()
}

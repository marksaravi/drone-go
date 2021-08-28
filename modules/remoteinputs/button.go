package remoteinputs

type button interface {
	Read() bool
}

type buttonInput struct {
	input button
	value bool
}

func (btn *buttonInput) read() (isChanged bool) {
	isChanged = false
	pv := btn.value
	btn.value = btn.input.Read()
	isChanged = btn.value != pv
	return
}

func (ri *remoteInputs) readStopButtons() (isChanged bool) {
	isLeftChanged := ri.buttonFrontLeft.read()
	return isLeftChanged
}

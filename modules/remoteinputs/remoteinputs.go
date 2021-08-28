package remoteinputs

type button interface {
	Read() bool
	Value() bool
}

type remoteInputs struct {
	buttonFrontLeft button
	stopped         bool
}

func NewRemoteInputs(buttonFrontLeft button) *remoteInputs {
	return &remoteInputs{
		buttonFrontLeft: buttonFrontLeft,
	}
}

func (ri *remoteInputs) RefreshInputs() (isStopChanged bool) {
	isStopChanged = ri.readStopButtons()
	return
}

func (ri *remoteInputs) readStopButtons() bool {
	var isChanged bool = false
	ri.buttonFrontLeft.Read()
	if ri.buttonFrontLeft.Value() {
		if !ri.stopped {
			isChanged = true
		}
		ri.stopped = true
	}
	return isChanged
}

func (ri *remoteInputs) IsStopped() bool {
	return ri.stopped
}

func (ri *remoteInputs) PrintData() {

}

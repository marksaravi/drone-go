package remoteinputs

import "fmt"

type button interface {
	Read() bool
}

type remoteInputs struct {
	buttonFrontLeft      button
	buttonFrontLeftValue bool
	stopped              bool
}

func NewRemoteInputs(buttonFrontLeft button) *remoteInputs {
	return &remoteInputs{
		buttonFrontLeft: buttonFrontLeft,
		stopped:         false,
	}
}

func (ri *remoteInputs) RefreshInputs() (isStopChanged bool) {
	isStopChanged = ri.readStopButtons()
	return
}

func (ri *remoteInputs) readStopButtons() bool {
	var isChanged bool = false
	ri.buttonFrontLeftValue = ri.buttonFrontLeft.Read()
	if ri.buttonFrontLeftValue {
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
	fmt.Printf("Button-Front-Left: %t", ri.buttonFrontLeftValue)
	fmt.Println()
}

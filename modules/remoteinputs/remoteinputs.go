package remoteinputs

import (
	buttonModule "github.com/MarkSaravi/drone-go/modules/button"
)

type RemoteInputsConfig struct {
	ButtonFrontEndGPIO string
}

type button interface {
	Read() bool
}

type remoteInputs struct {
	buttonFrontLeft button
}

func NewRemoteInputs(config RemoteInputsConfig) *remoteInputs {
	return &remoteInputs{
		buttonFrontLeft: buttonModule.NewButton(config.ButtonFrontEndGPIO),
	}
}

func (ri *remoteInputs) ReadInputs() {
	ri.buttonFrontLeft.Read()
}

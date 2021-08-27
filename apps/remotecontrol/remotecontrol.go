package remotecontrol

import "github.com/MarkSaravi/drone-go/modules/remoteinputs"

type RemoteConfig struct {
	ButtonFrontEndGPIO string
}

type inputs interface {
	ReadInputs()
}

type remoteControl struct {
	inputs inputs
}

func NewRemoteControl(config RemoteConfig) *remoteControl {
	inputs := remoteinputs.NewRemoteInputs(remoteinputs.RemoteInputsConfig{
		ButtonFrontEndGPIO: config.ButtonFrontEndGPIO,
	})
	return &remoteControl{
		inputs: inputs,
	}
}

func (rc *remoteControl) Start() {
	for {
		rc.inputs.ReadInputs()
	}
}

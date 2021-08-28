package remotecontrol

import (
	"time"
)

type inputs interface {
	ReadInputs()
	PrintData()
}

type remoteControl struct {
	inputs inputs
}

func NewRemoteControl(inputs inputs) *remoteControl {
	return &remoteControl{
		inputs: inputs,
	}
}

func (rc *remoteControl) Start() {
	for {
		rc.inputs.ReadInputs()
		rc.inputs.PrintData()
		time.Sleep(250 * time.Millisecond)
	}
}

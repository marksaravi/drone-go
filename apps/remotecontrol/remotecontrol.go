package remotecontrol

import (
	"time"
)

type inputs interface {
	RefreshInputs() (isStopChanged bool)
	IsStopped() bool
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
		rc.inputs.RefreshInputs()
		rc.inputs.PrintData()
		time.Sleep(250 * time.Millisecond)
	}
}

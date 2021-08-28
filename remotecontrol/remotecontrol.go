package remotecontrol

import (
	"time"

	"github.com/MarkSaravi/drone-go/models"
)

type inputs interface {
	ReadInputs() models.RemoteControlData
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

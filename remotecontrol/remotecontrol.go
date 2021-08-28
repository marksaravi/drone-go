package remotecontrol

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/models"
)

type button interface {
	Read() models.ButtonData
}

type remoteControl struct {
	btnFrontLeft button
	data         models.RemoteControlData
}

func NewRemoteControl(btnFrontLeft button) *remoteControl {
	return &remoteControl{
		btnFrontLeft: btnFrontLeft,
	}
}

func (rc *remoteControl) Start() {
	for {
		rc.read()
		rc.showData()
		time.Sleep(250 * time.Millisecond)
	}
}

func (rc *remoteControl) read() {
	rc.data = models.RemoteControlData{
		ButtonFrontLeft: rc.btnFrontLeft.Read(),
	}
}

func (rc *remoteControl) showData() {
	fmt.Println(rc.data)
}

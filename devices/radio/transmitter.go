package radio

import (
	"fmt"

	"github.com/marksaravi/drone-go/constants"
)

type radioTransmitterLink interface {
	PowerOn()
	ClearStatus()
	TransmitterOn()
	Transmit(data []byte) error
	IsTransmitFailed(update bool) bool
}

type radioTransmitter struct {
	radiolink radioTransmitterLink
}

func NewTransmitter(radiolink radioTransmitterLink) *radioTransmitter {
	return &radioTransmitter{
		radiolink: radiolink,
	}
}

func (t *radioTransmitter) On() {
	t.radiolink.TransmitterOn()
	t.radiolink.PowerOn()
}

func (t *radioTransmitter) Transmit(payload []byte) error {
	if len(payload) != int(constants.RADIO_PAYLOAD_SIZE) {
		return fmt.Errorf("radio: payload size is %d", len(payload))
	}
	if t.radiolink.IsTransmitFailed(true) {
		t.radiolink.ClearStatus()
	}
	return t.radiolink.Transmit(payload)
}

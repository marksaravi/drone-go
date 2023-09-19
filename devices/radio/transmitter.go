package radio

import (
	"fmt"
)

type radioTransmitterLink interface {
	PayloadSize() int
	PowerOn()
	PowerOff()
	ClearStatus()
	TransmitterOn()
	Transmit(data []byte) error
	IsTransmitFailed(update bool) bool
}

type radioTransmitter struct {
	radiolink radioTransmitterLink
}

func NewTransmitter(radiolink radioTransmitterLink, connectionTimeoutMs int) *radioTransmitter {
	return &radioTransmitter{
		radiolink: radiolink,
	}
}

func (t *radioTransmitter) Start() {
	t.radiolink.TransmitterOn()
	t.radiolink.PowerOn()
}

func (t *radioTransmitter) Transmit(payload []byte) error {
	if len(payload) != t.PayloadSize() {
		return fmt.Errorf("radio: payload size is %d", len(payload))
	}
	if t.radiolink.IsTransmitFailed(true) {
		t.radiolink.ClearStatus()
	}
	return t.radiolink.Transmit(payload)
}

func (t *radioTransmitter) PayloadSize() int {
	return t.radiolink.PayloadSize()
}

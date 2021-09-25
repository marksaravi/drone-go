package radioreceiver

import (
	"context"
	"testing"
	"time"

	"github.com/marksaravi/drone-go/models"
	"github.com/marksaravi/drone-go/utils"
)

const TIMEOUT_MS int = 250
const HEARBIT_PER_SEC int = 1
const RECEIVE_PER_SEC int = 50

type mockradio struct {
	mockdata        [][]byte
	receiveIndex    int
	isReceiverOn    bool
	isTransmitterOn bool
}

func (r *mockradio) Receive() ([]byte, bool) {
	i := r.receiveIndex
	r.receiveIndex++
	if i < len(r.mockdata) {
		return r.mockdata[i], true
	}
	return []byte{}, false
}

func (r *mockradio) ReceiverOn() {
	r.isReceiverOn = true
	r.isTransmitterOn = false
}

func (r *mockradio) TransmitterOn() {
	r.isReceiverOn = false
	r.isTransmitterOn = true
}

func (r *mockradio) Transmit(data []byte) error {
	return nil
}

func NewMockRadio(data [][]byte) *mockradio {
	return &mockradio{
		mockdata:     data,
		receiveIndex: 0,
	}
}

func TestReceiverConnected(t *testing.T) {
	radio := NewMockRadio([][]byte{utils.SerializeFlightCommand(models.FlightCommands{
		Id: 0,
	})})
	ctx, cancel := context.WithCancel(context.Background())
	receiver := NewRadioReceiver(ctx, radio, RECEIVE_PER_SEC, HEARBIT_PER_SEC, time.Millisecond*time.Duration(TIMEOUT_MS))
	var running bool = true
	time.AfterFunc(time.Duration(100)*time.Millisecond, func() {
		cancel()
	})
	var connected bool = false
	for running {
		select {
		case <-ctx.Done():
			running = false
		case <-receiver.command:
		case connected = <-receiver.connection:

		}
	}
	if !radio.isReceiverOn || radio.isTransmitterOn {
		t.Fatal("Receiver is not activated")
	}
	if !connected {
		t.Fatal("Receiver failed to connect")
	}
}

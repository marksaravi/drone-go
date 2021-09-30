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

type mockData struct {
	interval  time.Duration
	available bool
	data      [32]byte
}
type mockradio struct {
	startTime       time.Time
	tries           []time.Duration
	mockdata        []mockData
	receiveIndex    int
	isReceiverOn    bool
	isTransmitterOn bool
	cancel          context.CancelFunc
}

func (r *mockradio) Receive() ([32]byte, bool) {
	if r.receiveIndex == len(r.mockdata) {
		r.cancel()
		return [32]byte{}, false
	}
	r.tries = append(r.tries, time.Since(r.startTime))
	if time.Since(r.startTime) < r.mockdata[r.receiveIndex].interval {

		return [32]byte{}, false
	}
	r.startTime = time.Now()
	data := r.mockdata[r.receiveIndex].data
	available := r.mockdata[r.receiveIndex].available
	r.receiveIndex++
	return data, available
}

func (r *mockradio) ReceiverOn() {
	r.isReceiverOn = true
	r.isTransmitterOn = false
}

func (r *mockradio) TransmitterOn() {
	r.isReceiverOn = false
	r.isTransmitterOn = true
}

func (r *mockradio) Transmit(data [32]byte) error {
	return nil
}

func NewMockRadio(cancel context.CancelFunc, data []mockData) *mockradio {
	return &mockradio{
		mockdata:     data,
		receiveIndex: 0,
		cancel:       cancel,
		startTime:    time.Now(),
	}
}

func TestReceiverConnected(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	radio := NewMockRadio(cancel, []mockData{
		{
			data: utils.SerializeFlightCommand(models.FlightCommands{
				Id: 0,
			}),
			interval:  time.Second / time.Duration(RECEIVE_PER_SEC),
			available: true,
		},
		{
			data:      [32]byte{},
			interval:  time.Second / time.Duration(RECEIVE_PER_SEC),
			available: false,
		},
	})
	receiver := NewRadioReceiver(ctx, radio, RECEIVE_PER_SEC, HEARBIT_PER_SEC)
	var running bool = true
	var connected bool = false
	for running {
		select {
		case <-ctx.Done():
			running = false
		case <-receiver.Command:
		case connected = <-receiver.Connection:

		}
	}
	if !radio.isReceiverOn || radio.isTransmitterOn {
		t.Fatal("Receiver is not activated")
	}
	if !connected {
		t.Fatal("Receiver failed to connect")
	}
}

func TestReceiverTimeout(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	radio := NewMockRadio(cancel, []mockData{
		{
			data: utils.SerializeFlightCommand(models.FlightCommands{
				Id: 0,
			}),
			interval:  time.Second / time.Duration(RECEIVE_PER_SEC),
			available: true,
		},
		{
			data:      [32]byte{},
			interval:  time.Second/time.Duration(RECEIVE_PER_SEC) + time.Millisecond*time.Duration(TIMEOUT_MS),
			available: false,
		},
	})

	receiver := NewRadioReceiver(ctx, radio, RECEIVE_PER_SEC, HEARBIT_PER_SEC)
	var running bool = true
	var connected []bool = []bool{}
	for running {
		select {
		case <-ctx.Done():
			running = false
		case <-receiver.Command:
		case conn := <-receiver.Connection:
			t.Log(conn)
			connected = append(connected, conn)
		}
	}
	if !radio.isReceiverOn || radio.isTransmitterOn {
		t.Fatal("Receiver is not activated")
	}
	if len(connected) != 2 || !connected[0] || connected[1] {
		t.Fatal("Receiver failed to timeout", radio.tries, connected)
	}
}

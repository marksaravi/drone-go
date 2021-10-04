package radiotransmitter

import (
	"context"
	"testing"
	"time"
)

const HEARBIT_TIMEOUT_MS int = 250
const TRANSMIT_PER_SEC int = 20

type mockdata struct {
	interval  time.Duration
	available bool
	data      [32]byte
}

type mockradio struct {
	data            []mockdata
	receiveIndex    int
	cancel          context.CancelFunc
	startTime       time.Time
	isReceiverOn    bool
	isTransmitterOn bool
}

func (r *mockradio) Receive() ([32]byte, bool) {
	if r.receiveIndex == len(r.data) {
		r.cancel()
		return [32]byte{}, false
	}
	if time.Since(r.startTime) < r.data[r.receiveIndex].interval {
		return [32]byte{}, false
	}
	r.startTime = time.Now()
	data := r.data[r.receiveIndex].data
	available := r.data[r.receiveIndex].available
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

func NewMockRadio(cancel context.CancelFunc, data []mockdata) *mockradio {
	return &mockradio{
		data:         data,
		receiveIndex: 0,
		cancel:       cancel,
		startTime:    time.Now(),
	}
}

func TestTransmitterConnected(t *testing.T) {
}

func TestReceiverTimeout(t *testing.T) {
}

package radiotransmitter

import (
	"context"
	"sync"
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
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	radio := NewMockRadio(cancel, []mockdata{
		{
			data:      [32]byte{},
			interval:  time.Second,
			available: true,
		},
	})
	rt := NewRadioTransmitter(
		ctx,
		&wg,
		radio,
		TRANSMIT_PER_SEC,
		time.Millisecond*time.Duration(HEARBIT_TIMEOUT_MS),
	)
	var running bool = true
	var heartbeating bool = false
	for running {
		select {
		case <-ctx.Done():
			running = false
		case heartbeating = <-rt.DroneHeartBeat:
		}
	}
	if !radio.isReceiverOn || radio.isTransmitterOn {
		t.Fatal("Receiver is not activated")
	}
	if !heartbeating {
		t.Fatal("Receiver failed to get Heartbeat")
	}
}

func TestReceiverTimeout(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	radio := NewMockRadio(cancel, []mockdata{
		{
			data:      [32]byte{},
			interval:  time.Second,
			available: true,
		},
		{
			data:      [32]byte{},
			interval:  time.Millisecond * time.Duration(HEARBIT_TIMEOUT_MS),
			available: false,
		},
	})
	rt := NewRadioTransmitter(
		ctx,
		&wg,
		radio,
		TRANSMIT_PER_SEC,
		time.Millisecond*time.Duration(HEARBIT_TIMEOUT_MS),
	)
	var running bool = true
	var heartbeatings []bool = []bool{}
	for running {
		select {
		case <-ctx.Done():
			running = false
		case beat := <-rt.DroneHeartBeat:
			heartbeatings = append(heartbeatings, beat)
		}
	}
	if !radio.isReceiverOn || radio.isTransmitterOn {
		t.Fatal("Transmitter is not activated")
	}
	if len(heartbeatings) != 2 || !heartbeatings[0] || heartbeatings[1] {
		t.Fatal("Transmitter failed to get Heartbeat timeout", heartbeatings)
	}
}

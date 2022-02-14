package radio

import (
	"testing"
	"time"
)

const testtimeout = time.Second / 2

func TestConnectionStateConnected(t *testing.T) {
	prevState := DISCONNECTED
	connected := true
	lastConected := time.Now().Add(-time.Second)

	state, lastConnection := newConnectionState(connected, prevState, lastConected, testtimeout)
	if state != CONNECTED || lastConnection == lastConected {
		t.Fatalf("wanted CONNECTED, have %s\n", StateToString(state))
	}
}

func TestConnectionStateIdleAndDisconnected(t *testing.T) {
	prevState := IDLE
	connected := false
	lastConected := time.Now().Add(-time.Second)

	state, _ := newConnectionState(connected, prevState, lastConected, testtimeout)
	if state != WAITING_FOR_CONNECTION {
		t.Fatalf("wanted WAITING_FOR_CONNECTION, have %s\n", StateToString(state))
	}
}

func TestConnectionStateWaitingAndDisconnected(t *testing.T) {
	prevState := WAITING_FOR_CONNECTION
	connected := false
	lastConected := time.Now().Add(-time.Second)

	state, _ := newConnectionState(connected, prevState, lastConected, testtimeout)
	if state != WAITING_FOR_CONNECTION {
		t.Fatalf("wanted WAITING_FOR_CONNECTION, have %s\n", StateToString(state))
	}
}

func TestConnectionStateConnectedAndDisconnected(t *testing.T) {
	prevState := CONNECTED
	connected := false
	lastConected := time.Now().Add(-time.Second)

	state, _ := newConnectionState(connected, prevState, lastConected, testtimeout)
	if state != DISCONNECTED {
		t.Fatalf("wanted DISCONNECTED, have %s\n", StateToString(state))
	}
}

func TestConnectionStateDisconnectedAndDisconnected(t *testing.T) {
	prevState := DISCONNECTED
	connected := false
	lastConected := time.Now().Add(-time.Second)

	state, _ := newConnectionState(connected, prevState, lastConected, testtimeout)
	if state != DISCONNECTED {
		t.Fatalf("wanted DISCONNECTED, have %s\n", StateToString(state))
	}
}

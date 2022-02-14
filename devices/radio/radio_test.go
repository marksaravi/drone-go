package radio

import (
	"testing"
	"time"

	"github.com/marksaravi/drone-go/constants"
)

const testtimeout = time.Second / 2

func TestConnectionStateConnected(t *testing.T) {
	prevState := constants.DISCONNECTED
	connected := true
	lastConected := time.Now().Add(-time.Second)

	state, lastConnection := newConnectionState(connected, prevState, lastConected, testtimeout)
	if state != constants.CONNECTED || lastConnection == lastConected {
		t.Fatalf("wanted CONNECTED, have %s\n", StateToString(state))
	}
}

func TestConnectionStateIdleAndDisconnected(t *testing.T) {
	prevState := constants.IDLE
	connected := false
	lastConected := time.Now().Add(-time.Second)

	state, _ := newConnectionState(connected, prevState, lastConected, testtimeout)
	if state != constants.WAITING_FOR_CONNECTION {
		t.Fatalf("wanted WAITING_FOR_CONNECTION, have %s\n", StateToString(state))
	}
}

func TestConnectionStateWaitingAndDisconnected(t *testing.T) {
	prevState := constants.WAITING_FOR_CONNECTION
	connected := false
	lastConected := time.Now().Add(-time.Second)

	state, _ := newConnectionState(connected, prevState, lastConected, testtimeout)
	if state != constants.WAITING_FOR_CONNECTION {
		t.Fatalf("wanted WAITING_FOR_CONNECTION, have %s\n", StateToString(state))
	}
}

func TestConnectionStateConnectedAndDisconnected(t *testing.T) {
	prevState := constants.CONNECTED
	connected := false
	lastConected := time.Now().Add(-time.Second)

	state, _ := newConnectionState(connected, prevState, lastConected, testtimeout)
	if state != constants.DISCONNECTED {
		t.Fatalf("wanted DISCONNECTED, have %s\n", StateToString(state))
	}
}

func TestConnectionStateDisconnectedAndDisconnected(t *testing.T) {
	prevState := constants.DISCONNECTED
	connected := false
	lastConected := time.Now().Add(-time.Second)

	state, _ := newConnectionState(connected, prevState, lastConected, testtimeout)
	if state != constants.DISCONNECTED {
		t.Fatalf("wanted DISCONNECTED, have %s\n", StateToString(state))
	}
}

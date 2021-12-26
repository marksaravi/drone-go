package radio

import (
	"testing"
	"time"
)

func TestIdleToConnectByData(t *testing.T) {
	r := NewRadio(nil, 500)
	r.connectionState = IDLE
	go func() {
		<-r.connection
	}()
	r.setConnectionState(true, DATA_PAYLOAD)
	if r.connectionState != CONNECTED {
		t.Fatalf("Wanted CONNECTED, got %v", r.connectionState)
	}
}

func TestDisconneToConnectByData(t *testing.T) {
	r := NewRadio(nil, 500)
	r.connectionState = DISCONNECTED
	go func() {
		<-r.connection
	}()
	r.setConnectionState(true, DATA_PAYLOAD)
	if r.connectionState != CONNECTED {
		t.Fatalf("Wanted CONNECTED, got %v", r.connectionState)
	}
}

func TestDisconneToConnectByHeartBeat(t *testing.T) {
	r := NewRadio(nil, 500)
	r.connectionState = DISCONNECTED
	go func() {
		<-r.connection
	}()
	r.setConnectionState(true, HEARTBEAT_PAYLOAD)
	select {
	case <-r.connection:
	default:
	}
	if r.connectionState != CONNECTED {
		t.Fatalf("Wanted CONNECTED, got %v", r.connectionState)
	}
}

func TestConnectedToDisconnectByReceiverOff(t *testing.T) {
	r := NewRadio(nil, 500)
	r.connectionState = CONNECTED
	go func() {
		<-r.connection
	}()
	r.setConnectionState(true, RECEIVER_OFF_PAYLOAD)
	if r.connectionState != DISCONNECTED {
		t.Fatalf("Wanted DISCONNECTED, got %v", r.connectionState)
	}
}

func TestConnectedToLostByTimeout(t *testing.T) {
	timeoutMS := 500
	r := NewRadio(nil, timeoutMS)
	r.connectionState = CONNECTED
	r.lastReceivedHeartBeat = time.Now().Add(-time.Duration(timeoutMS * int(time.Millisecond)))
	go func() {
		<-r.connection
	}()
	r.setConnectionState(false, NO_PAYLOAD)
	if r.connectionState != LOST {
		t.Fatalf("Wanted LOST, got %v", r.connectionState)
	}
}

func TestLostToConnectByData(t *testing.T) {
	timeoutMS := 500
	r := NewRadio(nil, timeoutMS)
	r.connectionState = LOST
	r.lastReceivedHeartBeat = time.Now().Add(-time.Duration(timeoutMS * int(time.Millisecond)))
	go func() {
		<-r.connection
	}()
	r.setConnectionState(true, DATA_PAYLOAD)
	if r.connectionState != CONNECTED {
		t.Fatalf("Wanted CONNECTED, got %v", r.connectionState)
	}
}

func TestLostToDisconnectByData(t *testing.T) {
	timeoutMS := 500
	r := NewRadio(nil, timeoutMS)
	r.connectionState = LOST
	r.lastReceivedHeartBeat = time.Now().Add(-time.Duration(timeoutMS * int(time.Millisecond)))
	go func() {
		<-r.connection
	}()
	r.setConnectionState(true, RECEIVER_OFF_PAYLOAD)
	if r.connectionState != DISCONNECTED {
		t.Fatalf("Wanted DISCONNECTED, got %v", r.connectionState)
	}
}

func TestDisconnectedToDisconnectedByTimeout(t *testing.T) {
	timeoutMS := 500
	r := NewRadio(nil, timeoutMS)
	r.connectionState = DISCONNECTED
	r.lastReceivedHeartBeat = time.Now().Add(-time.Duration(timeoutMS * int(time.Millisecond)))
	go func() {
		<-r.connection
	}()
	r.setConnectionState(false, NO_PAYLOAD)
	if r.connectionState != DISCONNECTED {
		t.Fatalf("Wanted DISCONNECTED, got %v", r.connectionState)
	}
}

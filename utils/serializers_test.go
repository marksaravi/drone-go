package utils

import (
	"testing"

	"github.com/marksaravi/drone-go/models"
)

func TestFlighCommandsToByteArray(t *testing.T) {

	got := SerializeFlightCommand(models.FlightCommands{
		PayloadType:       27,
		Id:                133,
		Time:              1632550150486903000,
		Roll:              -3.23,
		Pitch:             4.17,
		Yaw:               -0.34,
		Throttle:          2.75,
		ButtonFrontLeft:   true,
		ButtonFrontRight:  false,
		ButtonTopLeft:     true,
		ButtonTopRight:    false,
		ButtonBottomLeft:  false,
		ButtonBottomRight: true,
	})
	want := [32]byte{27, 133, 0, 0, 0, 216, 72, 200, 85, 182, 251, 167, 22, 82, 184, 78, 192, 164, 112, 133, 64, 123, 20, 174, 190, 0, 0, 48, 64, 37, 0, 0}
	if !compareByteArrays(got[:], want[:]) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestFlighCommandsFromByteArray(t *testing.T) {
	got := DeserializeFlightCommand([32]byte{33, 133, 0, 0, 0, 216, 72, 200, 85, 182, 251, 167, 22, 82, 184, 78, 192, 164, 112, 133, 64, 123, 20, 174, 190, 0, 0, 48, 64, 37, 0, 0})
	want := models.FlightCommands{
		PayloadType:       33,
		Id:                133,
		Time:              1632550150486903000,
		Roll:              -3.23,
		Pitch:             4.17,
		Yaw:               -0.34,
		Throttle:          2.75,
		ButtonFrontLeft:   true,
		ButtonFrontRight:  false,
		ButtonTopLeft:     true,
		ButtonTopRight:    false,
		ButtonBottomLeft:  false,
		ButtonBottomRight: true,
	}
	if want != got {
		t.Errorf("got %v, want %v", got, want)
	}
}

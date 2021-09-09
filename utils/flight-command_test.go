package utils

import (
	"testing"

	"github.com/MarkSaravi/drone-go/models"
)

func TestFlighCommandsToByteArray(t *testing.T) {
	got := SerializeFlightCommand(models.FlightCommands{
		Id:                133,
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
	want := []byte{133, 0, 0, 0, 82, 184, 78, 192, 164, 112, 133, 64, 123, 20, 174, 190, 0, 0, 48, 64, 37, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	if !compareByteArrays(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestFlighCommandsFromByteArray(t *testing.T) {
	got := DeserializeFlightCommand([]byte{133, 0, 0, 0, 82, 184, 78, 192, 164, 112, 133, 64, 123, 20, 174, 190, 0, 0, 48, 64, 37, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	want := models.FlightCommands{
		Id:                133,
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

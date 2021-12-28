package utils

import (
	"testing"

	"github.com/marksaravi/drone-go/models"
)

func TestFlighCommandsToByteArray(t *testing.T) {

	got := SerializeFlightCommand(models.FlightCommands{
		PayloadType:       27,
		Id:                133,
		Roll:              12,
		Pitch:             13,
		Yaw:               112,
		Throttle:          250,
		ButtonFrontLeft:   true,
		ButtonFrontRight:  false,
		ButtonTopLeft:     true,
		ButtonTopRight:    false,
		ButtonBottomLeft:  false,
		ButtonBottomRight: true,
	})
	want := models.Payload{27, 133, 12, 13, 112, 250, 37, 0}
	if !compareByteArrays(got[:], want[:]) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestFlighCommandsFromByteArray(t *testing.T) {
	got := DeserializeFlightCommand(models.Payload{33, 133, 44, 45, 46, 47, 37, 0})
	want := models.FlightCommands{
		PayloadType:       33,
		Id:                133,
		Roll:              44,
		Pitch:             45,
		Yaw:               46,
		Throttle:          47,
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

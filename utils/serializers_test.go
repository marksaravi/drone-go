package utils

import (
	"testing"

	"github.com/marksaravi/drone-go/models"
)

func TestSerializeFlightCommand(t *testing.T) {

	got := SerializeFlightCommand(models.FlightCommands{
		PayloadType:       27,
		Roll:              369,
		Pitch:             815,
		Yaw:               519,
		Throttle:          1020,
		ButtonFrontLeft:   true,
		ButtonFrontRight:  false,
		ButtonTopLeft:     true,
		ButtonTopRight:    false,
		ButtonBottomLeft:  false,
		ButtonBottomRight: true,
	})
	want := models.Payload{27, 0, 37, 113, 47, 7, 252, 237}
	if !compareByteArrays(got[:], want[:]) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestDeserializeFlightCommand(t *testing.T) {
	got := DeserializeFlightCommand(models.Payload{33, 0, 37, 247, 247, 247, 247, 255})
	want := models.FlightCommands{
		PayloadType:       33,
		Roll:              1015,
		Pitch:             1015,
		Yaw:               1015,
		Throttle:          1015,
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

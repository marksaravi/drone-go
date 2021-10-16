package pidcontrol

import (
	"testing"

	"github.com/marksaravi/drone-go/models"
)

var conversions = analogToDigitalConversions{
	roll: analogToDigitalConversion{
		ratio:  1.5,
		offset: -2.5,
	},
	pitch: analogToDigitalConversion{
		ratio:  1.2,
		offset: -2.65,
	},
	yaw: analogToDigitalConversion{
		ratio:  1.4,
		offset: -2.3,
	},
	throttle: analogToDigitalConversion{
		ratio:  1.1,
		offset: 0.7,
	},
}

func TestCurrAndPrevRotaions(t *testing.T) {
	rotations1 := models.ImuRotations{
		Rotations: models.Rotations{
			Roll:  1,
			Pitch: 1,
			Yaw:   1,
		},
	}
	rotations2 := models.ImuRotations{
		Rotations: models.Rotations{
			Roll:  2,
			Pitch: 2,
			Yaw:   2,
		},
	}

	pid := pidControl{
		prevRotations: models.ImuRotations{
			Rotations: models.Rotations{
				Roll:  0,
				Pitch: 0,
				Yaw:   0,
			},
		},
		rotations: rotations1,
	}
	pid.ApplyRotations(rotations2)
	if pid.prevRotations != rotations1 || pid.rotations != rotations2 {
		t.Fatalf("Expectd to set rotations to %v  and prevRotations to %v but has %v and %v", rotations2, rotations1, pid.rotations, pid.prevRotations)
	}
}

func TestConversion(t *testing.T) {
	want := float64(3)
	got := convert(float32(1), analogToDigitalConversion{ratio: 2.3, offset: 0.7})
	if got != want {
		t.Fatalf("wanted %f, got %f", want, got)
	}
}

func TestFlightControlCommandToPIDCommand(t *testing.T) {
	want := pidCommands{
		roll:     -1,
		pitch:    -0.85,
		yaw:      4,
		throttle: 4,
	}
	got := flightControlCommandToPIDCommand(models.FlightCommands{
		Roll:     1,
		Pitch:    1.5,
		Yaw:      4.5,
		Throttle: 3,
	}, conversions)
	if got != want {
		t.Fatalf("wanted %v, got %v", want, got)
	}
}

func TestApplyFlightCommand(t *testing.T) {
	want := pidCommands{
		roll:     -1,
		pitch:    -0.85,
		yaw:      4,
		throttle: 4,
	}
	pid := pidControl{
		conversions: conversions,
	}
	pid.ApplyFlightCommands(models.FlightCommands{
		Roll:     1,
		Pitch:    1.5,
		Yaw:      4.5,
		Throttle: 3,
	})
	got := pid.commands
	if got != want {
		t.Fatalf("wanted %v, got %v", want, got)
	}
	if pid.commands.throttle != want.throttle {
		t.Fatalf("wanted throttle %f, got %f", want.throttle, pid.commands.throttle)
	}
}

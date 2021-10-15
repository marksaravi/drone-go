package pidcontrol

import (
	"testing"

	"github.com/marksaravi/drone-go/models"
)

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

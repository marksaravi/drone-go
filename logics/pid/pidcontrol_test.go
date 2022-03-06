package pid

import (
	"testing"
)

func TestCreate(t *testing.T) {
	got := NewPIDControls(
		PIDSettings{MaxOutputToMaxThrottleRatio: 0.5, PGain: 1.5, IGain: 2, DGain: 2.5},
		PIDSettings{MaxOutputToMaxThrottleRatio: 0.75, PGain: 2.5, IGain: 3, DGain: 3.5},
		PIDSettings{MaxOutputToMaxThrottleRatio: 0.2, PGain: 3.5, IGain: 4, DGain: 4.5},
		80,
		0.5,
		CalibrationSettings{
			Calibrating: "roll",
			Gain:        "d",
			PStep:       0.1,
			IStep:       0.2,
			DStep:       0.3,
		},
	)
	if got.rollPIDControl.pGain != 1.5 {
		t.Errorf("got pgain %f, want pgain %f", got.rollPIDControl.pGain, float64(1.5))
	}
	if got.rollPIDControl.maxOutput != 40 {
		t.Errorf("got max output %f, want max output %f", got.rollPIDControl.pGain, float64(40))
	}
	if got.rollPIDControl.maxI != 20 {
		t.Errorf("got pgain %f, want pgain %f", got.rollPIDControl.pGain, float64(1.5))
	}
	if got.pitchPIDControl.maxI != 30 {
		t.Errorf("got pgain %f, want pgain %f", got.rollPIDControl.pGain, float64(1.5))
	}
}

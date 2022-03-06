package pid

import (
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	got := NewPIDControls(
		PIDSettings{MaxOutputToMaxThrottleRatio: 0.5, PGain: 1.5, IGain: 2, DGain: 2.5},
		PIDSettings{MaxOutputToMaxThrottleRatio: 0.75, PGain: 2.5, IGain: 3, DGain: 3.5},
		PIDSettings{MaxOutputToMaxThrottleRatio: 0.2, PGain: 3.5, IGain: 4, DGain: 4.5},
		80,
		10,
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

func TestPIDGains(t *testing.T) {
	pidcontrols := NewPIDControls(
		PIDSettings{MaxOutputToMaxThrottleRatio: 0.75, PGain: 2.5, IGain: 0.5, DGain: 0.75},
		PIDSettings{MaxOutputToMaxThrottleRatio: 0.75, PGain: 2.5, IGain: 0.5, DGain: 0.75},
		PIDSettings{MaxOutputToMaxThrottleRatio: 0.75, PGain: 2.5, IGain: 0.5, DGain: 0.75},
		100,
		10,
		1,
		CalibrationSettings{
			Calibrating: "none",
			Gain:        "",
			PStep:       0.1,
			IStep:       0.2,
			DStep:       0.3,
		},
	)

	pValue := pidcontrols.rollPIDControl.getP(1.5)
	var wantP float64 = 3.75
	if pValue != wantP {
		t.Errorf("want roll P %f, got %f", wantP, pValue)
	}

	pidcontrols.rollPIDControl.iMemory = 5
	wantI := 5.00075
	iValue := pidcontrols.rollPIDControl.getI(1.5, time.Second/1000)
	if iValue != wantI {
		t.Errorf("want roll I %f, got %f", wantI, iValue)
	}

	wantD := float64(750)
	pidcontrols.rollPIDControl.previousInput = 1
	dValue := pidcontrols.rollPIDControl.getD(2, time.Second/1000)
	if dValue != wantD {
		t.Errorf("want roll D %f, got %f", wantD, dValue)
	}
}

func TestPIDLimits(t *testing.T) {
	pidcontrols := NewPIDControls(
		PIDSettings{MaxOutputToMaxThrottleRatio: 0.75, PGain: 10, IGain: 2, DGain: 3},
		PIDSettings{MaxOutputToMaxThrottleRatio: 0.75, PGain: 10, IGain: 2, DGain: 3},
		PIDSettings{MaxOutputToMaxThrottleRatio: 0.75, PGain: 10, IGain: 2, DGain: 3},
		100,
		10,
		1,
		CalibrationSettings{
			Calibrating: "none",
			Gain:        "",
			PStep:       0.1,
			IStep:       0.2,
			DStep:       0.3,
		},
	)

	pidcontrols.rollPIDControl.iMemory = pidcontrols.rollPIDControl.maxI
	wantI := pidcontrols.rollPIDControl.maxI
	iValue := pidcontrols.rollPIDControl.getI(1, time.Second/1000)
	if iValue != wantI {
		t.Errorf("want roll positive I %f, got %f", wantI, iValue)
	}

	pidcontrols.rollPIDControl.iMemory = -pidcontrols.rollPIDControl.maxI
	wantI = -pidcontrols.rollPIDControl.maxI
	iValue = pidcontrols.rollPIDControl.getI(-1, time.Second/1000)
	if iValue != wantI {
		t.Errorf("want roll negative I %f, got %f", wantI, iValue)
	}

	rollFeedback := pidcontrols.rollPIDControl.calcPIDFeedback(30, time.Second)
	var wantFeedback float64 = 75
	if rollFeedback != wantFeedback {
		t.Errorf("want roll feedback %f, got %f", wantFeedback, rollFeedback)
	}

	rollFeedback = pidcontrols.rollPIDControl.calcPIDFeedback(-30, time.Second)
	wantFeedback = 10
	if rollFeedback != wantFeedback {
		t.Errorf("want roll feedback %f, got %f", wantFeedback, rollFeedback)
	}
}

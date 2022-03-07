package pid

import (
	"testing"
	"time"
)

// func TestCreate(t *testing.T) {
// 	got := NewPIDControls(
// 		PIDSettings{PGain: 1.5, IGain: 2, DGain: 2.5},
// 		PIDSettings{PGain: 2.5, IGain: 3, DGain: 3.5},
// 		PIDSettings{PGain: 3.5, IGain: 4, DGain: 4.5},
// 		80,
// 		0.5,
// 		CalibrationSettings{
// 			Calibrating: "roll",
// 			Gain:        "d",
// 			PStep:       0.1,
// 			IStep:       0.2,
// 			DStep:       0.3,
// 		},
// 	)
// 	if got.rollPIDControl.pGain != 1.5 {
// 		t.Errorf("got pgain %f, want pgain %f", got.rollPIDControl.pGain, float64(1.5))
// 	}
// 	if got.rollPIDControl.maxI != 20 {
// 		t.Errorf("got pgain %f, want pgain %f", got.rollPIDControl.pGain, float64(1.5))
// 	}
// 	if got.pitchPIDControl.maxI != 30 {
// 		t.Errorf("got pgain %f, want pgain %f", got.rollPIDControl.pGain, float64(1.5))
// 	}
// }

func TestPIDGains(t *testing.T) {
	pidcontrols := NewPIDControls(
		PIDSettings{PGain: 2.5, IGain: 0.5, DGain: 0.75, MaxI: 10},
		PIDSettings{},
		PIDSettings{},
		CalibrationSettings{},
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
	pidcontrols.rollPIDControl.dPrevError = 1
	dValue := pidcontrols.rollPIDControl.getD(2, time.Second/1000)
	if dValue != wantD {
		t.Errorf("want roll D %f, got %f", wantD, dValue)
	}
}

func TestPIDLimits(t *testing.T) {
	pidcontrols := NewPIDControls(
		PIDSettings{PGain: 2, IGain: 3, DGain: 4, MaxI: 10},
		PIDSettings{},
		PIDSettings{},
		CalibrationSettings{},
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
}

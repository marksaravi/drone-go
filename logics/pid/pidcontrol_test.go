package pid

import (
	"testing"
	"time"
)

func TestPIDGains(t *testing.T) {
	pidcontrols := NewPIDControls(
		PIDSettings{PGain: 2.5, IGain: 0.5, DGain: 0.75, MaxI: 10},
		PIDSettings{},
		PIDSettings{},
		true,
		true,
		10,
		CalibrationSettings{},
	)

	pValue := pidcontrols.xPIDControl.getP(1.5)
	var wantP float64 = 3.75
	if pValue != wantP {
		t.Errorf("want roll P %f, got %f", wantP, pValue)
	}

	pidcontrols.xPIDControl.iMemory = 5
	wantI := 5.00075
	iValue := pidcontrols.xPIDControl.getI(1.5, time.Second/1000)
	if iValue != wantI {
		t.Errorf("want roll I %f, got %f", wantI, iValue)
	}

	wantD := float64(750)
	pidcontrols.xPIDControl.dPrevError = 1
	dValue := pidcontrols.xPIDControl.getD(2, time.Second/1000)
	if dValue != wantD {
		t.Errorf("want roll D %f, got %f", wantD, dValue)
	}
}

func TestPIDLimits(t *testing.T) {
	pidcontrols := NewPIDControls(
		PIDSettings{PGain: 2, IGain: 3, DGain: 4, MaxI: 10},
		PIDSettings{},
		PIDSettings{},
		true,
		true,
		10,
		CalibrationSettings{},
	)

	pidcontrols.xPIDControl.iMemory = pidcontrols.xPIDControl.maxI
	wantI := pidcontrols.xPIDControl.maxI
	iValue := pidcontrols.xPIDControl.getI(1, time.Second/1000)
	if iValue != wantI {
		t.Errorf("want roll positive I %f, got %f", wantI, iValue)
	}

	pidcontrols.xPIDControl.iMemory = -pidcontrols.xPIDControl.maxI
	wantI = -pidcontrols.xPIDControl.maxI
	iValue = pidcontrols.xPIDControl.getI(-1, time.Second/1000)
	if iValue != wantI {
		t.Errorf("want roll negative I %f, got %f", wantI, iValue)
	}
}

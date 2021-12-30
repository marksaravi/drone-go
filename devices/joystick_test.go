package devices

import "testing"

func TestJoystickCalcCoefficients(t *testing.T) {
	a, b := calcCoefficients(485, 1024)
	if a != -0.0001032840 || b != 1.1057628670 {
		t.Fatalf("failed with %15.10f, %15.10f", a, b)
	}
}

func TestJoystickCalcValue(t *testing.T) {
	var digitalMidValue uint16 = 485
	a, b := calcCoefficients(digitalMidValue, 1024)
	got := calcValue(digitalMidValue, a, b, 1)
	if got != 512 {
		t.Fatalf("failed with %d", got)
	}
}

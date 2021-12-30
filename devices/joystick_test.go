package devices

import "testing"

func TestJoystickCalcCoefficients(t *testing.T) {
	a, b := calcCoefficients(485, 1024)
	if a != -0.00002317 || b != 1.105762867 {
		t.Fatalf("failed with %15.10f, %15.10f", a, b)
	}
}

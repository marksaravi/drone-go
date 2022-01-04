package pid

import "testing"

func TestJoystickToPidValue(t *testing.T) {
	var digitalValue uint16 = 64
	pidcontrols := NewPIDControls(3200, 0, 0, 0, 16, 16, 16, 16, 1024)
	var want float64 = -7
	got := pidcontrols.joystickToPidValue(digitalValue, pidcontrols.rollPid.limit)
	if got != want {
		t.Fatalf("Wanted %3.2f, got %3.2f", want, got)
	}

}

func TestThrottleToPidThrottle(t *testing.T) {
	var digitalValue uint16 = 64
	pidcontrols := NewPIDControls(3200, 0, 0, 0, 16, 16, 16, 16, 1024)
	var want float64 = 1
	got := pidcontrols.throttleToPidThrottle(digitalValue)
	if got != want {
		t.Fatalf("Wanted %3.2f, got %3.2f", want, got)
	}
}

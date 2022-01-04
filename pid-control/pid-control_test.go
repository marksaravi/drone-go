package pidcontrol

import "testing"

func TestJoystickToPidValue(t *testing.T) {
	var digitalValue uint16 = 64
	pid := NewPIDControl(3200, 0, 0, 0, 16, 16, 16, 16, 1024)
	var want float64 = -7
	got := pid.joystickToPidValue(digitalValue, pid.maxRoll)
	if got != want {
		t.Fatalf("Wanted %3.2f, got %3.2f", want, got)
	}

}

func TestThrottleToPidThrottle(t *testing.T) {
	var digitalValue uint16 = 64
	pid := NewPIDControl(3200, 0, 0, 0, 16, 16, 16, 16, 1024)
	var want float64 = 1
	got := pid.throttleToPidThrottle(digitalValue)
	if got != want {
		t.Fatalf("Wanted %3.2f, got %3.2f", want, got)
	}

}

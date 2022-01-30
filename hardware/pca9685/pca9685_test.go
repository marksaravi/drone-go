package pca9685

import "testing"

func TestThrottleToPulseWidthPositive(t *testing.T) {
	var throttle float64 = MaxAllowedThrottle - 1
	want := MinPW + (MaxPW-MinPW)*throttle/100
	got := throttleToPulseWidth(throttle)
	if got != want {
		t.Fatalf("wanted %f, got %f", want, got)
	}
	throttle = 10
	want = MinPW + (MaxPW-MinPW)*0.1
	got = throttleToPulseWidth(throttle)
	if got != want {
		t.Fatalf("wanted %f, got %f", want, got)
	}
	throttle = MaxAllowedThrottle + 1
	want = MinPW + (MaxPW-MinPW)*MaxAllowedThrottle/100
	got = throttleToPulseWidth(throttle)
	if got != want {
		t.Fatalf("wanted %f, got %f", want, got)
	}
}

func TestThrottleToPulseWidthNegative(t *testing.T) {
	var throttle float64 = -50
	want := MinPW
	got := throttleToPulseWidth(throttle)
	if got != want {
		t.Fatalf("wanted %f, got %f", want, got)
	}
	throttle = -10
	got = throttleToPulseWidth(throttle)
	if got != want {
		t.Fatalf("wanted %f, got %f", want, got)
	}
}

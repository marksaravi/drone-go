package pca9685

import "testing"

func TestThrottleInRange(t *testing.T) {
	var throttle float64 = MaxAllowedThrottle - 1
	want := MinPW + (MaxPW-MinPW)*throttle/100
	got := throttleToPulseWidth(throttle)
	if got != want {
		t.Fatalf("wanted %f, got %f", want, got)
	}
}

func TestThrottleMoreThanMax(t *testing.T) {
	var throttle float64 = MaxAllowedThrottle + 1
	got := throttleToPulseWidth(throttle)
	if got != MaxAllowedPulseWidth {
		t.Fatalf("wanted %f, got %f", MaxAllowedPulseWidth, got)
	}
}

func TestThrottleEqualToMax(t *testing.T) {
	var throttle float64 = MaxAllowedThrottle
	got := throttleToPulseWidth(throttle)
	if got != MaxAllowedPulseWidth {
		t.Fatalf("wanted %f, got %f", MaxAllowedPulseWidth, got)
	}
}

func TestThrottleEqualToZero(t *testing.T) {
	var throttle float64 = 0
	got := throttleToPulseWidth(throttle)
	if got != MinPW {
		t.Fatalf("wanted %f, got %f", MinPW, got)
	}
}

func TestThrottleLessThanZero(t *testing.T) {
	var throttle float64 = -1
	got := throttleToPulseWidth(throttle)
	if got != MinPW {
		t.Fatalf("wanted %f, got %f", MinPW, got)
	}
}

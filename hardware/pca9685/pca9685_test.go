package pca9685

import "testing"

func TestThrottleInRange(t *testing.T) {
	var throttle float64 = MaxAllowedThrottle - 1
	got := limitThrottle(throttle, true)
	if got != throttle {
		t.Fatalf("wanted %f, got %f", throttle, got)
	}
}

func TestThrottleMoreThanMax(t *testing.T) {
	var throttle float64 = MaxAllowedThrottle + 1
	got := limitThrottle(throttle, true)
	if got != MaxAllowedThrottle {
		t.Fatalf("wanted %f, got %f", MaxAllowedThrottle, got)
	}
}

func TestThrottlessThanMin(t *testing.T) {
	var throttle float64 = MinOnThrottle - 1
	got := limitThrottle(throttle, true)
	if got != MinOnThrottle {
		t.Fatalf("wanted %f, got %f", MinOnThrottle, got)
	}
}

func TestThrottleWhenOff(t *testing.T) {
	var throttle float64 = MaxAllowedThrottle - 1
	got := limitThrottle(throttle, false)
	if got != 0 {
		t.Fatalf("wanted %f, got %f", 0.0, got)
	}
}

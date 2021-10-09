package utils

import (
	"context"
	"testing"
	"time"
)

func TestTicker(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	const TICK_PER_SEC int = 4000
	INTERVAL := int64(time.Second / time.Duration(TICK_PER_SEC))
	TOLERANCE := int64(float64(INTERVAL) * float64(0.05))
	MAX_INTERVAL := INTERVAL + TOLERANCE
	ticker := NewTicker(ctx, "1", TICK_PER_SEC)
	const NUM_OF_TICKS = 100
	ticks := make([]int64, NUM_OF_TICKS)
	for i := 0; i < NUM_OF_TICKS; i++ {
		ticks[i] = <-ticker
	}
	cancel()
	min := int64(time.Second)
	max := int64(time.Nanosecond)
	for i := 1; i < NUM_OF_TICKS; i++ {
		diff := ticks[i] - ticks[i-1]
		if diff < min {
			min = diff
		}
		if diff > max {
			max = diff
		}
	}
	if max > MAX_INTERVAL {
		t.Fatalf("Want %v-%v, Got %v-%v", INTERVAL, MAX_INTERVAL, min, max)
	}
}

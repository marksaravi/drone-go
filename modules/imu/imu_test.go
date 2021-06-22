package imu_test

import (
	"testing"

	"github.com/MarkSaravi/drone-go/modules/imu"
)

func TestInt16To2sComplement(t *testing.T) {
	const x int16 = -60
	got := imu.OffsetCorrection(5, 4, 0.9)
	const want float64 = 4.9
	if got != want {
		t.Errorf("got %f, want %f", got, want)
	}
}

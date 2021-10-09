package icm20948

import (
	"testing"
)

func TestHighLowBytesToPositiveInt16(t *testing.T) {
	got := towsComplementUint8ToInt16(64, 168)
	const want int16 = 16552
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func TestHighLowBytesToNegativeInt16(t *testing.T) {
	got := towsComplementUint8ToInt16(255, 125)
	const want int16 = -131
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

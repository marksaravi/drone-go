package icm20789_test

import (
	"testing"

	"github.com/marksaravi/drone-go/hardware/mems/icm20789"
)

func TestOffsetToHighAndLowBits(t *testing.T) {
	const offset uint16 = 3379
	expectedHigherBits := byte(26)
	expectedLowerBits := byte(102)
	higherBits, lowerBits := icm20789.OffsetToHL(offset)

	if higherBits != expectedHigherBits || lowerBits != expectedLowerBits {
		t.Errorf("Expected higherBits: %d, lowerBits: %d, got higherBits: %d, lowerBits: %d", expectedHigherBits, expectedLowerBits, higherBits, lowerBits)
	}

}

func TestHighAndLowBitsToOffset(t *testing.T) {
	const expectedOffset uint16 = 2236
	offset := icm20789.HLtoOffset(17, 120)

	if offset != expectedOffset {
		t.Errorf("Expected expectedOffset: %d, got offset: %d", expectedOffset, offset)
	}

}

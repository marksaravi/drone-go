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

func TestAccOffsetToLandHBytes(t *testing.T) {
	var offset uint16 = 31979
	var wantH uint8 = 249
	var wantL uint8 = 214
	gotH, gotL := accOffsetToHighandLowBytes(offset)
	if wantH != gotH || wantL != gotL {
		t.Errorf("got %d, %d, want %d, %d", gotH, gotL, wantH, wantL)
	}
}

package utils_test

import (
	"testing"

	"github.com/MarkSaravi/drone-go/utils"
)

/*
[254 132 255 196 64 168 255 241 0 37 0 7]
accX: -380, accY: -60, accZ: 16552
gyroX: -15, gyroY: 37, gyroZ: 7
254 132
255 196
64 168
*/

func TestInt16To2sComplement(t *testing.T) {
	const x int16 = -60
	got := utils.IntToTowsComplement(x)
	const want uint16 = 65476
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func TestInt16ToHighLowBytes(t *testing.T) {
	gotH, gotL := utils.IntToTowsComplementBytes(int16(-380))
	const wantH byte = 254
	const wantL byte = 132
	if gotH != wantH || gotL != wantL {
		t.Errorf("got %d %d, want %d %d", gotH, gotL, wantH, wantL)
	}
}

func TestHighLowBytesToInt16(t *testing.T) {
	got := utils.TowsComplementUint8ToInt16(64, 168)
	const want int16 = 16552
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

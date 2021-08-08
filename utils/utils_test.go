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

func TestHighLowBytesToInt16(t *testing.T) {
	got := utils.TowsComplementUint8ToInt16(64, 168)
	const want int16 = 16552
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

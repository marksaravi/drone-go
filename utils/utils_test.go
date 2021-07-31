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

func TestFloatArrayToByteArrayAndReverse(t *testing.T) {
	fa1 := []float32{366.34, -180.24, 0, -144.32, 22.22}
	ba := utils.FloatArrayToByteArray(fa1)
	fa2 := utils.ByteArrayToFloat32Array(ba)
	var equal bool = true
	for i := 0; i < len(fa1); i++ {
		if fa1[i] != fa2[i] {
			equal = false
		}
	}
	if !equal {
		t.Errorf("got %v, want %v", fa1, fa2)
	}
}

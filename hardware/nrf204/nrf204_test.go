package nrf204

import "testing"

func TestFloatArrayToByteArrayAndReverse(t *testing.T) {
	fa1 := []float32{366.34, -180.24, 0, -144.32, 22.22}
	ba := floatArrayToByteArray(fa1)
	fa2 := byteArrayToFloat32Array(ba)
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

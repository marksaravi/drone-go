package mpu_test

import (
	"reflect"
	"testing"

	"github.com/MarkSaravi/drone-go/modules/mpu"
	"github.com/MarkSaravi/drone-go/types"
)

func TestNew(t *testing.T) {
	data := mpu.New(4)
	got := len(data.GetBuffer())
	want := 4
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func TestPushToFront(t *testing.T) {
	data := mpu.New(4)
	for i := 1; i <= 4; i++ {
		data.PushToFront(types.XYZ{X: float64(i), Y: float64(i), Z: float64(i)})
	}
	want := []types.XYZ{
		{X: 4, Y: 4, Z: 4},
		{X: 3, Y: 3, Z: 3},
		{X: 2, Y: 2, Z: 2},
		{X: 1, Y: 1, Z: 1},
	}
	got := data.GetBuffer()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("got %v, want %v", got, want)
	}
}

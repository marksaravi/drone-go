package mpu_test

import (
	"reflect"
	"testing"

	"github.com/MarkSaravi/drone-go/modules/mpu"
	"github.com/MarkSaravi/drone-go/types"
)

func TestNew(t *testing.T) {
	data := mpu.NewSensorData(4)
	got := len(data.GetBuffer())
	want := 4
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func TestPushToFront(t *testing.T) {
	data := mpu.NewSensorData(4)
	want2 := []types.XYZ{
		{X: 2, Y: 2, Z: 2},
		{X: 1, Y: 1, Z: 1},
		{X: 0, Y: 0, Z: 0},
		{X: 0, Y: 0, Z: 0},
	}
	want4 := []types.XYZ{
		{X: 4, Y: 4, Z: 4},
		{X: 3, Y: 3, Z: 3},
		{X: 2, Y: 2, Z: 2},
		{X: 1, Y: 1, Z: 1},
	}
	data.PushToFront(types.XYZ{X: 1, Y: 1, Z: 1})
	data.PushToFront(types.XYZ{X: 2, Y: 2, Z: 2})
	if !reflect.DeepEqual(want2, data.GetBuffer()) {
		t.Errorf("got %v, want %v", data.GetBuffer(), want4)
	}
	data.PushToFront(types.XYZ{X: 3, Y: 3, Z: 3})
	data.PushToFront(types.XYZ{X: 4, Y: 4, Z: 4})
	if !reflect.DeepEqual(want4, data.GetBuffer()) {
		t.Errorf("got %v, want %v", data.GetBuffer(), want4)
	}
}
